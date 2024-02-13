package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pe-Gomes/simple-bank-go/util"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../../")

	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("Could not create connection pool with DB: ", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
