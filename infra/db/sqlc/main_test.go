package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:changeme@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("Could not create connection pool with DB: ", err)
	}

	testQueries = New(connPool)
	os.Exit(m.Run())
}
