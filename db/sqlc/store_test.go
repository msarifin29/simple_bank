package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := CreateAccount(t)
	account2 := CreateAccount(t)
	fmt.Println("before >> ", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entity
		fromEntity := result.FromEntry
		require.NotEmpty(t, fromEntity)
		require.Equal(t, account1.ID, fromEntity.AccountID)
		require.Equal(t, -amount, fromEntity.Amount)
		require.NotZero(t, fromEntity.CreatedAt)
		require.NotZero(t, fromEntity.ID)

		_, err = testStore.GetEntry(context.Background(), fromEntity.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.CreatedAt)
		require.NotZero(t, toEntry.ID)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		_, err = testStore.GetAccount(context.Background(), fromAccount.ID)
		require.NoError(t, err)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		_, err = testStore.GetAccount(context.Background(), toAccount.ID)
		require.NoError(t, err)

		// check balance
		fmt.Println(`>> tx :`, fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	// check final update balance
	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(`>> after :`, updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}
func TestTransferTxDeadLock(t *testing.T) {
	account1 := CreateAccount(t)
	account2 := CreateAccount(t)
	fmt.Println("before >> ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccount := account1.ID
		toAccount := account2.ID

		if i%2 == 1 {
			fromAccount = account2.ID
			toAccount = account1.ID
		}

		go func() {

			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount,
				ToAccountID:   toAccount,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check final update balance
	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(`>> after :`, updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
