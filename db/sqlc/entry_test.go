package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/vadym-98/simple_bank/util"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T) Entry {
	account1 := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account1.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account1 := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		arg := CreateEntryParams{
			account1.ID,
			util.RandomMoney(),
		}
		e, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, e)
	}

	arg := ListEntriesParams{
		account1.ID,
		5,
		5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, e := range entries {
		require.NotEmpty(t, e)
	}
}
