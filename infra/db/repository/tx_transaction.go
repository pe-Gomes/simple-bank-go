package db

import "context"

type TransactionTxResult struct {
	Transaction Transaction `json:"transaction"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
	FromEntry   Entry       `json:"from_entry"`
	ToEntry     Entry       `json:"to_entry"`
}

func (store *SQLStore) TransactionTx(ctx context.Context, arg CreateTransactionParams) (TransactionTxResult, error) {
	var result TransactionTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accoundID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accoundID1,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
