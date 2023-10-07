package db

import (
	"context"
	"gochat/src/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) CreateUserTxResult {
	arg := CreateUserTxParams{
		Username: util.RandomString(6),
		Password: util.RandomString(6),
		Role:     "admin",
	}

	res, err := testStore.CreateUserTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, res.User)

	require.Equal(t, arg.Role, res.User.Role)
	require.Equal(t, arg.Username, res.User.Username)
	require.NotEqual(t, arg.Password, res.User.HashedPassword)
	require.False(t, res.User.Deleted)

	return res
}

func TestCreateUserTx(t *testing.T) {
	res := createRandomUser(t)

	err := testStore.DeleteUserTx(context.Background(), res.User.ID)
	require.NoError(t, err)
}
