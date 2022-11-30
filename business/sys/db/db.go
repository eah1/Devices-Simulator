// Package db provides support for access the database.
package db

import (
	"device-simulator/app/config"
	"time"

	errors2 "device-simulator/business/sys/errors"
	goSentry "github.com/getsentry/sentry-go"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

// Config is the required properties to use the database.
type Config struct {
	DBPostgres   string
	MaxIdleConns int
	MaxOpenConns int
	ShowSQL      bool
}

// NewConfigDB create a Config Data Base.
func NewConfigDB(config config.Config) Config {
	configDB := new(Config)
	configDB.DBPostgres = config.DBPostgres
	configDB.MaxIdleConns = config.DBMaxIdleConns
	configDB.MaxOpenConns = config.DBMaxOpenConns
	configDB.ShowSQL = config.DBLogger

	return *configDB
}

// Open knows how to open a DB connection based on the configuration.
func Open(cfg Config, retries int, log *zap.SugaredLogger) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("pgx", cfg.DBPostgres)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new engine")
	}

	i := 1

	err = engine.Ping()
	for err != nil && i < retries {
		log.Warn("Failed to connect to database (%d)", retries)
		time.Sleep(time.Duration(500*i) * time.Millisecond)

		i++

		engine, _ = xorm.NewEngine("pgx", cfg.DBPostgres)
		err = engine.Ping()
	}

	if err != nil {
		return nil, errors.Wrap(err, "unable to connection database")
	}

	engine.SetMaxIdleConns(cfg.MaxIdleConns)
	engine.SetMaxOpenConns(cfg.MaxOpenConns)

	// show log query sql.
	engine.ShowSQL(cfg.ShowSQL)

	return engine, nil
}

func TranslatePsqlError(log *zap.SugaredLogger, err error) error {
	goSentry.CaptureException(err)

	switch e := err.(type) {
	case *pgconn.PgError:
		log.Errorw("Postgresql error", "service",
			"POSTGRESQL | DB", "constraints", e.ConstraintName, "error", err.Error())

		switch e.Code {
		case "22P02":
			// Invalid input syntax
			return errors2.ErrElementRequest
		case "23503":
			// Foreign key violation
			return errors2.ErrElementRequest
		case "23505":
			// Unique violation
			return errors2.ErrElementDuplicated
		case "42703":
			// Undefined column
			return errors2.ErrElementRequest
		default:
			return errors2.ErrPsqlDefault
		}
	default:
		return errors2.ErrPsqlDefault
	}
}
