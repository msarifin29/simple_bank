package db

import (
	"context"
	"testing"

	"github.com/msarifin29/simple_bank/util"
	"github.com/stretchr/testify/assert"
)

func CreateAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testStore.CreateAccount(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)
	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	CreateAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	assert.Nil(t, err)
	assert.NotNil(t, account2)
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	err := testStore.UpdateAccount(context.Background(), arg)
	assert.Nil(t, err)
}
func TestListAccount(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotNil(t, accounts)
	for _, account := range accounts {
		assert.NoError(t, err)
		assert.NotEmpty(t, account)

	}
}
func TestDeleteccount(t *testing.T) {
	err := testStore.DeleteAccount(context.Background(), 1)
	assert.Nil(t, err)
}
