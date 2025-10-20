package database

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
)

// IDB base interface can be used on behalf of both bun.DB and bun.Tx
type IDB interface {
	bun.IDB
}

type DB struct {
	*bun.DB
}

type Tx struct {
	*bun.Tx
}

func NewDB(cfg *config.Config, logger logging.Logger) *DB {
	return &DB{
		DB: connect(cfg, logger),
	}
}

func connect(cfg *config.Config, logger logging.Logger) *bun.DB {
	conf, err := pgx.ParseConfig(cfg.DB.GetDSN())
	if err != nil {
		logger.Error("Failed to parse dsn")
	}
	sqlDB := sql.OpenDB(stdlib.GetConnector(*conf))

	db := bun.NewDB(sqlDB, pgdialect.New())
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	if !cfg.IsProdEnv() {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to database: %v", err)
	}

	logger.Info("Successfully connected to db")

	return db
}
