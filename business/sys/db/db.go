// Package db provides support for access the database.
package db

import (
	"device-simulator/app/config"
	mycDBErrors "device-simulator/business/db/errors"
	mycErrors "device-simulator/business/sys/errors"
	"errors"
	"fmt"
	"time"

	goSentry "github.com/getsentry/sentry-go"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

const timeSleep = 500

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
		return nil, fmt.Errorf("db.open.NewEngine %w", err)
	}

	i := 1

	err = engine.Ping()
	for err != nil && i < retries {
		log.Warn("Failed to connect to database (%d)", retries)
		time.Sleep(time.Duration(timeSleep*i) * time.Millisecond)

		i++

		engine, _ = xorm.NewEngine("pgx", cfg.DBPostgres)
		err = engine.Ping()
	}

	if err != nil {
		return nil, fmt.Errorf("db.open.NewEngine %w", err)
	}

	engine.SetMaxIdleConns(cfg.MaxIdleConns)
	engine.SetMaxOpenConns(cfg.MaxOpenConns)

	// show log query sql.
	engine.ShowSQL(cfg.ShowSQL)

	return engine, nil
}

// PsqlError parse error PgError to custom error.
func PsqlError(log *zap.SugaredLogger, err error) error {
	var pgError *pgconn.PgError

	if errors.As(err, &pgError) {
		return &mycDBErrors.PsqlError{
			CodeSQL:        pgError.Code,
			TableName:      pgError.TableName,
			ConstraintName: pgError.ConstraintName,
			Err:            pgError.Message,
		}
	}

	log.Error(err)
	goSentry.CaptureException(err)

	return fmt.Errorf("db.PsqlError: %w", mycErrors.ErrDB)
}
