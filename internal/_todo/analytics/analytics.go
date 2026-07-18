package analytics

import (
	"context"
	"errors"
	"regexp"
	"sync"
	"time"

	"github.com/MindHunter86/eyesonly/internal/dynamic"
	"github.com/MindHunter86/eyesonly/internal/stats"
	"github.com/MindHunter86/eyesonly/internal/utils"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
)

type CallbackFunc func(lockdown bool)

type Analytics struct {
	db *analyticsDB

	filter *regexp.Regexp

	lm          sync.RWMutex
	analyzeint  *atomic.Duration
	analyzewnd  *atomic.Duration
	lockdowndur *atomic.Duration
	lockdownfrc *atomic.Bool

	analyzetrgrd  *atomic.Bool
	analyzettime  *atomic.Time
	analyzeerrcnt *atomic.Uint32

	trgrBootHold  *atomic.Duration
	trgrRatioWarn *atomic.Float64
	trgrRatioErr  *atomic.Float64

	cm        sync.RWMutex
	callbacks []CallbackFunc

	log *zerolog.Logger
	rdy *atomic.Bool
	sts *stats.Stats

	cfgupd chan struct{}
}

const defaultMetricFilter = "fiber\\.requests"

func NewAnalytics(c context.Context) (_ *Analytics, e error) {
	cli := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)

	a := &Analytics{
		log: utils.ContextValueExtract[*zerolog.Logger](c, utils.CtxZeroLogger),
		rdy: &atomic.Bool{},

		callbacks: make([]CallbackFunc, 0, 32),

		analyzetrgrd:  atomic.NewBool(false),
		analyzettime:  atomic.NewTime(time.Time{}),
		analyzeerrcnt: atomic.NewUint32(0),

		cfgupd: make(chan struct{}),

		// initializion via dynamic configuration below
		analyzeint:  atomic.NewDuration(0),
		analyzewnd:  atomic.NewDuration(0),
		lockdowndur: atomic.NewDuration(0),
		lockdownfrc: atomic.NewBool(false),

		trgrRatioWarn: atomic.NewFloat64(0),
		trgrRatioErr:  atomic.NewFloat64(0),
		trgrBootHold:  atomic.NewDuration(0),
	}

	// initial configuration load and callback registration
	a.onConfigurationUpdated(c)
	dynamic.RegisterCallback(a.onConfigurationUpdated)

	if a.analyzeint.Load()%(5*time.Second) != 0 {
		panic("application intervals must be multiples of 5 seconds (arg analytics-analyze-interval invalid)")
	}

	var expr string
	if expr = cli.String("analytics-metric-filter"); expr == "" {
		expr = defaultMetricFilter
	}

	if a.filter, e = regexp.Compile(expr); e != nil {
		return
	}

	a.db, e = newAnalyticsDB(c)
	return a, e
}

// Method for stats.StatsWriter interface for handling metrics
func (m *Analytics) WriteOne(t *time.Time, k string, v uint64) (int64, error) {
	if !m.rdy.Load() {
		m.log.Warn().Msg("analytics Write() skipped due to non-ready module state")
		return -1, nil
	}

	if !m.filter.MatchString(k) {
		return -1, nil
	}

	return m.db.writeMetrics(t, k, v, m.sts)
}

func (m *Analytics) onServiceTicker1sec(c context.Context) (e error) { // skipcq: GO-R1005 - here is no cyclomatic complexity
	tick := utils.ContextValueExtract[uint64](c, utils.CtxTickerTick)
	if tick%uint64(m.analyzeint.Load().Seconds()) != 0 {
		return nil
	} else if m.lockdownfrc.Load() {
		return nil
	}

	if !m.lm.TryLock() {
		m.sts.WriteIncMetric(stats.IMAnltLoopLockedError)
		if !m.analyzetrgrd.Load() {
			return errors.New("could not start loop event because of mutex lock")
		}
		return errors.New("could not start loop event because of mutex lock (lockdown mode)")
	}
	defer m.lm.Unlock()

	m.sts.WriteIncMetric(stats.IMAnltLoopCount)

	lap := utils.ContextValueExtract[time.Time](c, utils.CtxTickerLap)
	fromtime, totime := lap.Add(m.analyzewnd.Load()*-1), lap

	// recover last values before lockdown mode
	if ft := m.analyzettime.Load(); !ft.IsZero() {
		fromtime = ft
	}

	var trgrd bool
	if _, trgrd, e = m.analyze(&fromtime, &totime); e != nil {
		// wait 30 seconds (by default) and then discard lockdown mode until database recovers
		if m.analyzetrgrd.Load() && m.analyzeerrcnt.Load() >= 3 {
			m.forEachCallback(false)
			m.analyzetrgrd.Swap(false)
			m.analyzettime.Store(time.Time{})
		}

		if m.analyzeerrcnt.Load() < 3 {
			m.analyzeerrcnt.Add(1)
		}
		m.log.Warn().Msg(utils.ExtraErrorWrapper(e, "extracting analytics data from database").Error())
		return
	}

	// reset database errors
	if cnt := m.analyzeerrcnt.Swap(0); cnt != 0 {
		m.log.Info().Msg("analytics database errors were discarded")
	}

	// ignore all lockdown triggers on service bootstrap
	if tick < uint64(m.trgrBootHold.Load().Seconds()) {
		m.log.Trace().Msg("lockdown processing will be skipped due to boot holder")
		return nil
	}

	switch {
	// continue last lockdown
	case trgrd && m.analyzetrgrd.Load():
		m.log.Info().Msgf("lockdown mode was prolonged for %.2f seconds", m.lockdowndur.Load().Seconds())
		m.isCallbackInterrupted(c, time.NewTimer(m.lockdowndur.Load()))

	// start lockdown
	case trgrd && !m.analyzetrgrd.Load():
		m.sts.WriteIncMetric(stats.IMAnltLockdownEnable)
		m.log.Info().Msgf("lockdown mode enabled for %.2f seconds; all protectors were called",
			m.lockdowndur.Load().Seconds())

		m.forEachCallback(true)
		m.analyzetrgrd.Swap(true)
		m.analyzettime.Store(fromtime)
		m.isCallbackInterrupted(c, time.NewTimer(m.lockdowndur.Load()))

	// stop last lockdown
	case !trgrd && m.analyzetrgrd.Load():
		m.sts.WriteIncMetric(stats.IMAnltLockdownDisable)
		m.log.Info().Msg("lockdown mode was disabled, all protectors were called")

		m.forEachCallback(false)
		m.analyzetrgrd.Swap(false)
		m.analyzettime.Store(time.Time{})
	}

	return
}

func (m *Analytics) onConfigurationUpdated(c context.Context) {
	// interrupt ticker callback in lockdown case
	select {
	case m.cfgupd <- struct{}{}:
	default:
		m.log.Warn().Msg("analytics: non-blocking send to cfgupd channel failed: channel is full")
	}

	// block ticker callback
	m.lm.Lock()
	defer m.lm.Unlock()

	// pre-update old value save

	// local struct configuration update
	m.analyzeint.Store(dynamic.Duration(c, "analytics-analyze-interval"))
	m.analyzewnd.Store(dynamic.Duration(c, "analytics-analyze-window"))
	m.lockdowndur.Store(dynamic.Duration(c, "analytics-lockdown-time"))

	m.trgrBootHold.Store(dynamic.Duration(c, "analytics-trigger-boothold"))
	m.trgrRatioWarn.Store(dynamic.Float64(c, "analytics-trigger-ratio-warn"))
	m.trgrRatioErr.Store(dynamic.Float64(c, "analytics-trigger-ratio-err"))

	ldold := m.lockdownfrc.Swap(dynamic.Bool(c, "analytics-lockdown-force"))

	// ignore post-update on initial configuration setup
	if !m.rdy.Load() {
		return
	}

	// post-update trigger calls
	if ldcurr := m.lockdownfrc.Load(); ldold != ldcurr {
		if ldcurr {
			m.log.Info().Msg("lockdown mode manually enabled; all protectors were called")
			m.forEachCallback(true)
		} else {
			if m.analyzetrgrd.Load() {
				m.analyzetrgrd.Swap(false)
				m.analyzettime.Store(time.Time{})
			}
			m.forEachCallback(false)
			m.log.Info().Msg("lockdown mode was manually disabled; all protectors were called")
		}
	}
}

func (m *Analytics) analyze(from, to *time.Time) (_ *time.Time, triggered bool, e error) {
	now := time.Now()
	if from == nil || from.IsZero() {
		from = &now
	}

	var res []*sqlAnalyzeRow
	if res, e = m.db.getStatsDifferences(from, to, m.sts); e != nil {
		return
	}

	for _, r := range res {
		ratio := r.ratio.Float64

		m.sts.WriteAvgMetric(stats.AMAnltStatReqCurr, uint64(r.curr.Float64))
		m.sts.WriteAvgMetric(stats.AMAnltStatReqDiff, uint64(r.diff.Float64))
		m.sts.WriteAvgMetricFloat(stats.AMAnltStatReqRatio, ratio, 1000) // ratio is round(3) float

		if ratio >= m.trgrRatioWarn.Load() && ratio < m.trgrRatioErr.Load() {
			m.log.Warn().Msgf("analytics event - warning because of high ratio %.3f", ratio)
			m.log.Trace().Msgf("analytics results: time now %s, time prev %s",
				r.timenow.String, r.timeprev.String)
			m.log.Info().Msgf("analytics results: curr rps - %.2f, diff with prev - %.2f, ratio - %.3f",
				r.curr.Float64, r.diff.Float64, ratio)
		} else if ratio >= m.trgrRatioErr.Load() {
			m.log.Error().Msgf("analytics event - error because of too high ratio %.3f", ratio)
			m.log.Trace().Msgf("analytics results: time now %s, time prev %s",
				r.timenow.String, r.timeprev.String)
			m.log.Info().Msgf("analytics results: curr rps - %.2f, diff with prev - %.2f, ratio - %.3f",
				r.curr.Float64, r.diff.Float64, ratio)
			triggered = true
		} else if m.log.GetLevel() <= zerolog.DebugLevel {
			m.log.Trace().Msgf("analytics results: time now %s, time prev %s",
				r.timenow.String, r.timeprev.String)
			m.log.Trace().Msgf("analytics results: curr rps - %.2f, diff with prev - %.2f, ratio - %.3f",
				r.curr.Float64, r.diff.Float64, ratio)
		}
	}

	return &now, triggered, nil
}

func (m *Analytics) isCallbackInterrupted(c context.Context, t *time.Timer) {
	defer t.Stop()

	select {
	case <-c.Done():
		m.log.Trace().Msg("program abort caught, interrupt callback Timers")
	case <-t.C:
		m.log.Trace().Msg("ticker finished, interrupt callback")
	case <-m.cfgupd:
		m.log.Trace().Msg("config update called, interrupt callback Timers")
	}
}
