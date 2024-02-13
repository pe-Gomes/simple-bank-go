package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "go.uber.org/mock/mockgen/model"
)

type Store interface {
	Querier
	TransactionTx(ctx context.Context, arg CreateTransactionParams) (TransactionTxResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPoll *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPoll,
		Queries:  New(connPoll),
	}
}
