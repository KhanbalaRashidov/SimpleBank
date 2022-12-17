package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/KhanbalaRashidov/SimpleBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomMoney(),
	}

	entry, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.FromAccountID, entry.FromAccountID)
	require.Equal(t, arg.ToAccountID, entry.ToAccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	entryCreated := createRandomTransfer(t)
	entryQueried, err := testQueries.GetTransfer(context.Background(), entryCreated.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entryQueried)

	require.Equal(t, entryCreated.ID, entryQueried.ID)
	require.Equal(t, entryCreated.FromAccountID, entryQueried.FromAccountID)
	require.Equal(t, entryCreated.ToAccountID, entryQueried.ToAccountID)
	require.Equal(t, entryCreated.Amount, entryQueried.Amount)
	require.WithinDuration(t, entryCreated.CreatedAt, entryQueried.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	entry1 := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     entry1.ID,
		Amount: util.RandomMoney(),
	}

	entry2, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)
}

func TestDeleteTransfer(t *testing.T) {
	entry1 := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), entry1.ID)

	require.NoError(t, err)

	entry2, err := testQueries.GetTransfer(context.Background(), entry1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
