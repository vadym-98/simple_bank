package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/vadym-98/simple_bank/util"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T) Transfer {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)

	arg := CreateTransferParams{
		a1.ID,
		a2.ID,
		util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	t1 := createRandomTransfer(t)

	t2, err := testQueries.GetTransfer(context.Background(), t1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, t2)

	require.Equal(t, t1.ID, t2.ID)
	require.Equal(t, t1.FromAccountID, t2.FromAccountID)
	require.Equal(t, t1.ToAccountID, t2.ToAccountID)
	require.Equal(t, t1.Amount, t2.Amount)
	require.WithinDuration(t, t1.CreatedAt, t2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)
	a3 := createRandomAccount(t)
	fromAccounts := []Account{a1, a2, a3}

	b1 := createRandomAccount(t)
	b2 := createRandomAccount(t)
	b3 := createRandomAccount(t)
	toAccounts := []Account{b1, b2, b3}

	for _, a := range fromAccounts {
		for _, b := range toAccounts {
			arg := CreateTransferParams{
				a.ID,
				b.ID,
				util.RandomMoney(),
			}
			transfer, err := testQueries.CreateTransfer(context.Background(), arg)
			require.NoError(t, err)
			require.NotEmpty(t, transfer)
		}
	}

	for _, a := range fromAccounts {
		for _, b := range toAccounts {
			arg := ListTransfersParams{
				a.ID,
				b.ID,
				5,
				0,
			}
			transfers, err := testQueries.ListTransfers(context.Background(), arg)
			require.NoError(t, err)
			require.Len(t, transfers, 5)

			for _, tr := range transfers {
				require.NotEmpty(t, tr)
			}
		}
	}
}
