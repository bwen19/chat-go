package db

import (
	"context"
	"fmt"
	"gochat/src/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	DeleteUserTx(ctx context.Context, userID int64) error
}

// SqlStore provides all functions to execute SQL queries and transactions
type SqlStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(config *utils.Config) (Store, error) {
	err := runDatabaseMigration(config)
	if err != nil {
		return nil, err
	}

	connPool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool")
	}

	store := &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
	return store, nil
}

func runDatabaseMigration(config *utils.Config) error {
	migration, err := migrate.New(config.MigrationUrl, config.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up")
	}
	return nil
}
