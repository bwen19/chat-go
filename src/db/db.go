package db

import (
	"context"
	"errors"
	"gochat/src/db/rdb"
	"gochat/src/db/sqlc"
	"gochat/src/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// dbStore provides all functions to execute queries
type dbStore struct {
	connPool *pgxpool.Pool
	*sqlc.Queries
	*rdb.Caches
}

// New creates a new app store
func New(config *util.Config) (Store, error) {
	err := runDatabaseMigration(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, config.DatabaseUrl)
	if err != nil {
		return nil, errors.New("cannot create pgx pool")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       0,
	})

	store := &dbStore{
		connPool: connPool,
		Queries:  sqlc.New(connPool),
		Caches:   rdb.New(client),
	}

	if err = store.setupAdministrator(ctx); err != nil {
		return nil, errors.New("cannot create admin account")
	}

	return store, nil
}

// Add admin if not exists
func (s *dbStore) setupAdministrator(ctx context.Context) error {
	if _, err := s.RetrieveUserByName(ctx, "admin"); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			_, err := s.CreateUser(ctx, &CreateUserParams{
				Username: "admin",
				Password: "123456",
				Role:     RoleAdmin,
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

// ExecTx executes a function within a database transaction
func (s *dbStore) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.WithTx(tx)
	if err = fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
