package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Execute the transaction (Tx)
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{DeferrableMode: pgx.Deferrable})

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rb_err := tx.Rollback(ctx); rb_err != nil {
			return fmt.Errorf("Transaction error: %v. Rollback error: %v", err, rb_err)
		}
		return err
	}

	return tx.Commit(ctx)
}
