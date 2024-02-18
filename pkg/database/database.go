// Package database provides a database connection.
package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"strings"

	// Import the mysql database driver.
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

var (
	ErrDBNotFound        = errors.New("data not found")
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
)

// Config for a database connection.
type Config struct {
	// Type of database.
	//
	// e.g. "postgres", "mysql", etc.
	Type string

	// Credentials.
	User     string
	Password string

	// Host and port.
	Host string
	Port int

	// Name of the database.
	Name string

	// MaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	MaxIdleConns int

	// MaxOpenConns sets the maximum number of open connections to the database.
	MaxOpenConns int

	// DisableTLS disables TLS for the connection.
	DisableTLS bool
}

// DBResults to store database operation results
type DBResults struct {
	LastInsertID int64
	AffectedRows int64
}

// Open a connection to a db.
func Open(cfg Config) (*sqlx.DB, error) {
	cs := connectionString(cfg)
	log.Printf("Connecting to database: %s", cs)

	db, err := sqlx.Open(cfg.Type, cs)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("database: %w", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	log.Printf("Connected to driver: %s", db.DriverName())
	return db, nil
}

// connectionString returns a connection string for the given config.
func connectionString(cfg Config) string {
	if strings.ToLower(cfg.Type) != "postgres" {
		return ""
	}

	return pgConnectionString(cfg)
}

// pgConnectionString returns a connection string for Postgres.
// func pgConnectionString(cfg Config) string {
// 	sslMode := "require"
// 	if cfg.DisableTLS {
// 		sslMode = "disable"
// 	}

// 	return fmt.Sprintf(
// 		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s timezone=utc",
// 		cfg.User,
// 		cfg.Password,
// 		cfg.Host,
// 		cfg.Port,
// 		cfg.Name,
// 		sslMode,
// 	)
// }

// func pgConnectionString(cfg Config) string {
// 	sslMode := "require"
// 	if cfg.DisableTLS {
// 		sslMode = "disable"
// 	}

// 	q := make(url.Values)
// 	q.Set("sslmode", sslMode)
// 	q.Set("timezone", "utc")

// 	if cfg.Port > 0 {
// 		cfg.Host = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
// 	}

// 	u := url.URL{
// 		User:     url.UserPassword(cfg.User, cfg.Password),
// 		Host:     cfg.Host,
// 		Path:     cfg.Name,
// 		RawQuery: q.Encode(),
// 	}

// 	return u.String()
// }

func pgConnectionString(cfg Config) string {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s timezone=utc",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		sslMode,
	)
}

type Transactor interface {
	Beginx() (*sqlx.Tx, error)
}

func WithinTran(ctx context.Context, log *slog.Logger, db Transactor, fn func(sqlx.ExtContext) error) error {
	tx, err := db.Beginx()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to begin transaction: %v", err))
		return err
	}
	// Mark to the defer function a rollback is required.
	mustRollback := true

	// Set up a defer function for rolling back the transaction. If
	// mustRollback is true it means the call to fn failed, and we
	// need to roll back the transaction.
	defer func() {
		if mustRollback {
			if err := tx.Rollback(); err != nil {
				log.Error("unable to rollback db transaction", "ERROR", err)
			}
		}
	}()

	// Execute the code inside the transaction. If the function
	// fails, return the error and the defer function will roll back.
	if err := fn(tx); err != nil {
		return fmt.Errorf("exec db transaction: %w", err)
	}

	// Disarm the deferred rollback.
	mustRollback = false

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit db transaction: %w", err)
	}

	return nil
}

func NamedExecContext(ctx context.Context, log *slog.Logger, db sqlx.ExtContext, query string, data interface{}) (DBResults, error) {
	var dbres DBResults
	res, err := sqlx.NamedExecContext(ctx, db, query, data)
	if err != nil {
		return DBResults{}, err
	}

	// lid, err := res.LastInsertId()
	// if err != nil {
	// 	return DBResults{}, err
	// }
	// dbres.LastInsertID = lid

	ra, err := res.RowsAffected()
	if err != nil {
		return DBResults{}, err
	}
	dbres.AffectedRows = ra

	if val, err := res.RowsAffected(); err != nil {
		dbres.AffectedRows = val
	}

	return dbres, err

}

func NamedQuerySlice(ctx context.Context, log *slog.Logger, db sqlx.ExtContext, query string, data interface{}, dest interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	rows, err := sqlx.NamedQueryContext(ctx, db, query, data)
	if err != nil {
		return err
	}
	defer rows.Close() //nolint:all

	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil && !strings.Contains(err.Error(), "unsupported Scan, storing driver.Value type <nil> into type *json.RawMessage") {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}

// NamedQueryStruct is a helper function for executing queries that return a
// single value to be unmarshalled into a struct type.
func NamedQueryStruct(ctx context.Context, log *slog.Logger, db sqlx.ExtContext, query string, data interface{}, dest interface{}) error {
	rows, err := sqlx.NamedQueryContext(ctx, db, query, data)
	if err != nil {
		return err
	}
	defer rows.Close() //nolint:all

	if !rows.Next() {
		return ErrDBNotFound
	}

	if err := rows.StructScan(dest); err != nil && !strings.Contains(err.Error(), "unsupported Scan, storing driver.Value type <nil> into type *json.RawMessage") {
		return err
	}

	return nil
}
