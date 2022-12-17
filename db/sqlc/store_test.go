package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	store := NewStore(testDb)

	testAccount1 := createRandomAccount(t)
	testAccount2 := createRandomAccount(t)

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < 5; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: testAccount1.ID,
				ToAccountID:   testAccount2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < 5; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, testAccount1.ID, transfer.FromAccountID)
		require.Equal(t, testAccount2.ID, transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, testAccount1.ID, fromEntry.AccountID)
		require.Equal(t, amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, testAccount2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, testAccount1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, testAccount2.ID, toAccount.ID)

	}

}
