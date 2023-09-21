package db

import (
	"context"
	"gochat/src/utils"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal("failed to create pgx pool:", err)
	}

	testStore = &SqlStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
	os.Exit(m.Run())
}
