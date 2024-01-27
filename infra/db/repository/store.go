package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPoll *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPoll,
		Queries:  New(connPoll),
	}
}

// Execute the transaction (Tx)
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	queries := New(tx)

	err = fn(queries)

	if err != nil {
		if rb_err := tx.Rollback(ctx); rb_err != nil {
			return fmt.Errorf("Transaction error: %v. Rollback error: %v", err, rb_err)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransactionTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransactionTxResult struct {
	Transfer    Transaction `json:"transfer"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
	FromEntry   Entry       `json:"from_entry"`
	ToEntry     Entry       `json:"to_entry"`
}

// Perform a whole bank transaction wihtin one DB transaction.
func (store *Store) TransferTx(ctx context.Context, args TransactionTxParams) (
	TransactionTxResult, error,
) {
	var result TransactionTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create a new bank Transaction
		result.Transfer, err = q.CreateTransaction(ctx, CreateTransactionParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})

		if err != nil {
			return err
		}

		// Create a new Entry for whom is SENDING the transaction
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})

		if err != nil {
			return err
		}

		// Create a new Entry for whom is RECEIVING the transaction
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})

		if err != nil {
			return err
		}

		// TODO: update account's balance

		return nil
	})

	return result, err
}
