package config

import (
	"context"
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	Pool     *pgxpool.Pool
	DBSource string
	once     sync.Once
}

// NewPostgresClient creates a new PostgresClient instance and returns a DB pool.
// It uses sync.Once to ensure the DB pool is only initialized once.
func (p *PostgresClient) NewPostgresClient(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	p.once.Do(func() {
		p.Pool, err = pgxpool.New(ctx, p.DBSource)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return p.Pool, err
}

// PingDB pings the database server to check the connection.
// This method is a convenience wrapper around pgxpool.Pool.Ping(ctx).
func (p *PostgresClient) PingDB(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}

// RunDBMigration runs database migrations using the provided migration URL.
func (p *PostgresClient) RunDBMigration(migrationURL string) error {
	migration, err := migrate.New(migrationURL, p.DBSource)
	if err != nil {
		log.Println("cannot create new migrate instance")
		return err
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("failed to run migrate up")
		return err
	}

	log.Println("DB migrated successfully")
	return nil
}

// SetupDatabase creates a new PostgresClient instance, initializes the database connection pool,
// pings the database to check the connection, and runs any pending database migrations.
func SetupDatabase(ctx context.Context, cfg Config, migrationURL string) (*pgxpool.Pool, error) {
	postgresClient := PostgresClient{DBSource: cfg.DSN}

	dbPool, err := postgresClient.NewPostgresClient(ctx)
	if err != nil {
		return nil, err
	}

	if err = postgresClient.PingDB(ctx); err != nil {
		return nil, err
	}

	if err = postgresClient.RunDBMigration(migrationURL); err != nil {
		return nil, err
	}

	return dbPool, nil
}
