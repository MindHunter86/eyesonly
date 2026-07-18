package analytics

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/MindHunter86/eyesonly/internal/stats"
	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	_ "github.com/marcboeker/go-duckdb/v2"
)

type analyticsDB struct {
	*sql.DB

	dbpath   string
	tblstats string

	log *zerolog.Logger
	ctx context.Context
}

type sqlAnalyzeRow struct {
	dummyValue
	timenow, timeprev sql.NullString
	curr              sql.NullFloat64
	diff, ratio       sql.NullFloat64
}

type dummyValue struct{}

func (dummyValue) Scan(any) error {
	return nil
}

var errNullValReceived = errors.New("NULL value detected by scanner")

// !!! TODO CLI VALUES
// TODO DATABASE defer Close()
func newAnalyticsDB(c context.Context) (_ *analyticsDB, e error) {
	cli := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)

	adb := &analyticsDB{
		ctx: c,
		log: utils.ContextValueExtract[*zerolog.Logger](c, utils.CtxZeroLogger),

		dbpath:   cli.String("analytics-database-path"),
		tblstats: cli.String("analytics-database-table-stats"),
	}

	if adb.DB, e = sql.Open(sqlDatabaseDriver, adb.dbpath+sqlDatabaseParams); e != nil {
		return
	}

	if cli.Bool("analytics-database-reinit") {
		if adb.DB, e = adb.reinitDatabase(cli.Bool("analytics-database-i-am-sure")); e != nil {
			return
		}
	}

	if e = adb.Ping(); e != nil {
		adb.Close()
		return
	}

	return adb, adb.initDatabaseSchema() // skipcq: GO-W4006 - code minify
}

func (m *analyticsDB) reinitDatabase(confirm bool) (_ *sql.DB, e error) {
	if !confirm {
		// skipcq: SCC-ST1005 - humanized message required
		return nil, errors.New(`you have used analytics-database-reinit, which will completely destroy your data;
		add --analytics-database-i-am-sure if you really want to continue data wiping.`)
	}

	if e = m.Close(); e != nil {
		m.log.Warn().Msg(utils.ExtraErrorWrapper(e, "closing database").Error())
	}

	if e = os.Remove(m.dbpath); e != nil {
		return nil, utils.ExtraErrorWrapper(e, "removing database file with os.Remove()")
	}

	return sql.Open(sqlDatabaseDriver, m.dbpath+sqlDatabaseParams)
}

// ? migrations
func (m *analyticsDB) initDatabaseSchema() (e error) {
	// create tables
	tblschema := "(time TIMESTAMP, metric VARCHAR, value INTEGER)"
	if _, e = m.ExecContext(m.ctx, "CREATE TABLE IF NOT EXISTS "+m.tblstats+tblschema); e != nil {
		m.Close()
		return utils.ExtraErrorWrapper(e, "creating table stats in analytics database")
	}

	// create macroses
	for idx, schm := range dbSchemaMacros {
		if _, e = m.ExecContext(m.ctx, schm); e != nil {
			m.Close()
			return utils.ExtraErrorWrapper(e, "creating macros %d in analytics database", idx)
		}
	}

	return
}

func (m *analyticsDB) writeMetrics(t *time.Time, name string, val uint64, sts *stats.Stats) (_ int64, e error) {
	sbuf := utils.AcquireStringsBuffer()
	defer utils.ReleaseStringsBuffer(sbuf)

	sbuf.WriteString("INSERT INTO ")
	sbuf.WriteString(m.tblstats)
	sbuf.WriteString(" VALUES (")

	tm := t.UTC().UnixMicro()
	sbuf.WriteString("make_timestamp(")
	sbuf.WriteString(strconv.FormatInt(tm, 10))
	sbuf.WriteByte(')')
	sbuf.WriteByte(',')

	sbuf.WriteByte('\'')
	sbuf.WriteString(name)
	sbuf.WriteByte('\'')
	sbuf.WriteByte(',')

	sbuf.WriteString(strconv.FormatUint(val, 10))
	sbuf.WriteString(")")

	// !!
	// todo : use smth without result allocation
	sts.WriteIncMetric(stats.IMAnltDBRequests)
	_, e = m.ExecContext(m.ctx, sbuf.String())
	if e != nil {
		sts.WriteIncMetric(stats.IMAnltDBReqErrors)
	}

	return int64(sbuf.Len()), e
}

// SQL request like
// select * from get_1m_requested_differences('2025-10-13 08:40:04.507072', 'fiber.requests');
func (m *analyticsDB) getStatsDifferences(fromtime, totime *time.Time, sts *stats.Stats) (_ []*sqlAnalyzeRow, e error) {
	sbuf := utils.AcquireStringsBuffer()
	defer utils.ReleaseStringsBuffer(sbuf)

	sbuf.WriteString("select * from get_1m_requested_differences(")

	sbuf.WriteString("make_timestamp(")
	sbuf.WriteString(strconv.FormatInt(fromtime.UTC().Truncate(time.Second).UnixMicro(), 10))
	sbuf.WriteByte(')')
	sbuf.WriteByte(',')

	if totime != nil {
		sbuf.WriteString("make_timestamp(")
		sbuf.WriteString(strconv.FormatInt(totime.UTC().Truncate(time.Second).UnixMicro(), 10))
		sbuf.WriteByte(')')
		sbuf.WriteByte(',')
	} else {
		sbuf.WriteString("NULL")
		sbuf.WriteByte(',')
	}

	sbuf.WriteString("'fiber.requests')")

	m.log.Trace().Msgf("duckdb query execution: %s", sbuf.String())

	started := time.Now()
	sts.WriteIncMetric(stats.IMAnltDBRequests)

	var rows *sql.Rows
	if rows, e = m.QueryContext(m.ctx, sbuf.String()); e != nil {
		sts.WriteIncMetric(stats.IMAnltDBReqErrors)
		return
	}
	defer rows.Close()

	// elapsed; db request exec time
	sts.WriteAvgMetric(stats.AMAnltDBReqTime,
		uint64(time.Since(started).Nanoseconds()))

	var result []*sqlAnalyzeRow
	for rows.Next() {
		// acquire
		res := &sqlAnalyzeRow{}

		e = rows.Scan(&res.timenow, &res.timeprev, &res.curr, res.dummyValue, &res.diff, &res.ratio)

		if e != nil {
			m.log.Warn().Msg("BUG: something wrong in rows.Scan - " + e.Error())
			return
		}

		// dummy check, all NULLs are filtered in utils.go sql macros
		if !res.curr.Valid || !res.diff.Valid || !res.ratio.Valid {
			m.log.Warn().Msg(utils.ExtraErrorWrapper(errNullValReceived, "check received row for NULL").Error())
			continue
		}

		result = append(result, res)
	}

	return result, nil
}
