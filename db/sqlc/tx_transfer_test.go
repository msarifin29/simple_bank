package db

import (
	"context"
	"testing"

	"github.com/msarifin29/simple_bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotEmpty(t, transfer)
	assert.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	assert.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	assert.Equal(t, arg.Amount, transfer.Amount)
	assert.NotZero(t, transfer.ID)
	assert.NotZero(t, transfer.CreatedAt)
	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := CreateAccount(t)
	account2 := CreateAccount(t)
	CreateRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := CreateAccount(t)
	account2 := CreateAccount(t)
	transfer1 := CreateRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, transfer1)
	assert.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	assert.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	assert.Equal(t, transfer1.Amount, transfer2.Amount)
	assert.Equal(t, transfer1.ID, transfer2.ID)
}

func TestListTransfer(t *testing.T) {
	account1 := CreateAccount(t)
	account2 := CreateAccount(t)

	for i := 0; i < 5; i++ {
		CreateRandomTransfer(t, account1, account2)
		CreateRandomTransfer(t, account2, account1)
	}
	arg :=
		ListTransfersParams{
			FromAccountID: account1.ID,
			ToAccountID:   account1.ID,
			Limit:         5,
			Offset:        5,
		}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)

	for _, tf := range transfers {
		assert.NotEmpty(t, tf)
		require.True(t, tf.FromAccountID == account1.ID || tf.ToAccountID == account1.ID)
	}
}
