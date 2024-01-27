package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransactionTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// Run 'n' cuncurrent transfer transactions with a certain amount
	n := 3
	amount := int64(10)

	// Create channels to handle curcurrent transaction results and/or errors
	errs := make(chan error)
	results := make(chan TransactionTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransactionTx(context.Background(), TransactionTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()

		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)

			result := <-results
			require.NotEmpty(t, result)

			transaction := result.Transaction
			require.NotEmpty(t, transaction)
			require.Equal(t, account1.ID, transaction.FromAccountID)
			require.Equal(t, account2.ID, transaction.ToAccountID)
			require.Equal(t, amount, transaction.Amount)
			require.NotZero(t, transaction.ID)
			require.NotZero(t, transaction.CreatedAt)

			_, err = testStore.GetTransaction(context.Background(), transaction.ID)
			require.NoError(t, err)

			fromEntry := result.FromEntry
			require.NotEmpty(t, fromEntry)
			require.Equal(t, account1.ID, fromEntry.AccountID)
			require.Equal(t, -amount, fromEntry.Amount)
			require.NotZero(t, fromEntry.ID)
			require.NotZero(t, fromEntry.CreatedAt)

			_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
			require.NoError(t, err)

			toEntry := result.ToEntry
			require.NotEmpty(t, toEntry)
			require.Equal(t, account2.ID, toEntry.AccountID)
			require.Equal(t, amount, toEntry.Amount)
			require.NotZero(t, toEntry.ID)
			require.NotZero(t, toEntry.CreatedAt)

			_, err = testStore.GetEntry(context.Background(), toEntry.ID)
			require.NoError(t, err)

			// TODO: check account's balance
		}
	}
}
