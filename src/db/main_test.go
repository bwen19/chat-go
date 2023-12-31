package db

import (
	"context"
	"gochat/src/db/sqlc"
	"gochat/src/util"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal("failed to create pgx pool: ", err)
	}
	defer connPool.Close()

	testStore = &dbStore{
		connPool: connPool,
		Queries:  sqlc.New(connPool),
	}
	os.Exit(m.Run())
}
