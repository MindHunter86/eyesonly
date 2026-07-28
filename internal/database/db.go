package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MindHunter86/eyesonly/internal/utils"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
)

func Open(c context.Context) (_ *sql.DB, e error) {
	cli := utils.ContextValueExtract[*cli.Context](c, utils.CtxCliContext)

	driver := cli.String("database-driver")
	if driver == "sqlite" {
		driver = "sqlite3"
	}

	dsn := cli.String("database-dsn")
	if driver == "sqlite" {
		dsn = fmt.Sprintf("%s?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000", dsn)
	}

	var database *sql.DB
	if database, e = sql.Open(driver, dsn); e != nil {
		return
	}

	if driver == "sqlite" {
		database.SetMaxOpenConns(1)
		database.SetMaxIdleConns(1)
		database.SetConnMaxLifetime(0)
	} else {
		database.SetMaxOpenConns(25)
		database.SetMaxIdleConns(10)
		database.SetConnMaxLifetime(30 * time.Minute)
	}

	if e = database.Ping(); e != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping db: %w", e)
	}

	return database, nil
}
