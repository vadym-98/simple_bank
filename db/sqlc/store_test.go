package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)
	fmt.Println(">>> before:", a1.Balance, a2.Balance)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	results := make(chan TransferTxResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				a1.ID,
				a2.ID,
				amount,
			})

			results <- result
			errs <- err
		}()
	}

	existed := make(map[int]bool)

	// check results
	for i := 0; i < n; i++ {
		r := <-results
		e := <-errs

		require.NoError(t, e)
		require.NotEmpty(t, r)

		// check transfer
		transfer := r.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, a1.ID, transfer.FromAccountID)
		require.Equal(t, a2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := r.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, a1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := r.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, a2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := r.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, a1.ID, fromAccount.ID)

		toAccount := r.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, a2.ID, toAccount.ID)

		fmt.Println(">>> tx:", fromAccount.Balance, toAccount.Balance)

		// check accounts' balance
		diff1 := a1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - a2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)

	fmt.Println(">>> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, a1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, a2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)
	fmt.Println(">>> before:", a1.Balance, a2.Balance)

	// run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := a1.ID
		toAccountID := a2.ID

		if i%2 == 1 {
			fromAccountID = a2.ID
			toAccountID = a1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				fromAccountID,
				toAccountID,
				amount,
			})

			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		e := <-errs
		require.NoError(t, e)

	}

	// check the final updated balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)

	fmt.Println(">>> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, a1.Balance, updatedAccount1.Balance)
	require.Equal(t, a2.Balance, updatedAccount2.Balance)
}
