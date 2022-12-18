package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	store := NewStore(testDb)

	testAccount1 := createRandomAccount(t)
	testAccount2 := createRandomAccount(t)
	fmt.Println("before balance>>", testAccount1.Balance, testAccount2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: testAccount1.ID,
				ToAccountID:   testAccount2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	//check result
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, testAccount1.ID, transfer.FromAccountID)
		require.Equal(t, testAccount2.ID, transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entry
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

		//check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, testAccount1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, testAccount2.ID, toAccount.ID)

		//check balance
		fmt.Println("after tx>>", fromAccount.Balance, toAccount.Balance)
		diff1 := testAccount1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - testAccount2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= 5)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	//check update balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), testAccount1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), testAccount2.ID)
	require.NoError(t, err)

	fmt.Println("after balance>>", testAccount1.Balance, testAccount2.Balance)

	require.Equal(t, testAccount1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, testAccount2.Balance+int64(n)*amount, updateAccount2.Balance)

}
