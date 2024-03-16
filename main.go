package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pe-Gomes/simple-bank-go/api"
	db "github.com/pe-Gomes/simple-bank-go/infra/db/repository"
	"github.com/pe-Gomes/simple-bank-go/util"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("Could not create connection pool with DB: ", err)
	}

	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	
	if err != nil {
		log.Fatal("Could not create server: ", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Could not start server:", err)
	}
}
