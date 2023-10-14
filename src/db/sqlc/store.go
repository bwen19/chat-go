package db

import (
	"context"
	"errors"
	"gochat/src/util"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	ExecTx
	Close()
}

// SqlStore provides all functions to execute SQL queries and transactions
type SqlStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(config *util.Config) (Store, error) {
	err := runDatabaseMigration(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, config.DatabaseUrl)
	if err != nil {
		return nil, errors.New("failed to create pgx pool")
	}

	store := &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}

	if err = store.createAdmin(ctx); err != nil {
		return nil, errors.New("failed to create admin account")
	}

	return store, nil
}

func runDatabaseMigration(config *util.Config) error {
	migration, err := migrate.New(config.MigrationUrl, config.DatabaseUrl)
	if err != nil {
		return errors.New("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.New("failed to run migrate up")
	}
	return nil
}

func (s *SqlStore) createAdmin(ctx context.Context) error {
	if _, err := s.GetUserByName(ctx, "admin"); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err := s.CreateUserTx(ctx, CreateUserTxParams{
				Username: "admin",
				Password: "123456",
				Role:     "admin",
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (s *SqlStore) Close() {
	s.connPool.Close()
}
