package db

import (
	"context"
	"testing"

	"github.com/msarifin29/simple_bank/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "password",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.Time.IsZero())
	assert.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T) {
	CreateUser(t)
}

func GetUser(t *testing.T) {
	user1 := CreateUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	assert.Nil(t, err)
	assert.NotNil(t, user2)
	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)
}

func TestGetUser(t *testing.T) {
	GetUser(t)
}
