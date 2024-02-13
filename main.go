package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pe-Gomes/simple-bank-go/api"
	db "github.com/pe-Gomes/simple-bank-go/infra/db/repository"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:changeme@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	connPool, err := pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("Could not create connection pool with DB: ", err)
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Could not start server:", err)
	}
}
