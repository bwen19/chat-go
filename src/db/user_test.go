package db

import (
	"context"
	"gochat/src/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *UserInfo {
	arg := &CreateUserParams{
		Username: util.RandomString(6),
		Password: util.RandomString(6),
		Role:     "admin",
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Role, user.Role)
	require.Equal(t, arg.Username, user.Username)
	require.False(t, user.Deleted)

	return user
}

func TestCreateUserTx(t *testing.T) {
	user := createRandomUser(t)

	err := testStore.RemoveUser(context.Background(), user.ID)
	require.NoError(t, err)
}

// func TestXxx(t *testing.T) {
// 	oldUser := createRandomUser(t)

// 	newName := util.RandomString(8)
// 	newNick := util.RandomNumString(8)
// 	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
// 		ID: oldUser.ID,
// 		Username: pgtype.Text{
// 			String: newName,
// 			Valid:  true,
// 		},
// 		Nickname: pgtype.Text{
// 			String: newNick,
// 			Valid:  true,
// 		},
// 	})

// 	require.NoError(t, err)
// 	require.Equal(t, oldUser.ID, updatedUser.ID)
// 	require.Equal(t, newName, updatedUser.Username)
// 	require.Equal(t, newNick, updatedUser.Nickname)
// 	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
// }

// func TestListUsers(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomUser(t)
// 	}

// 	users, err := testStore.ListUsers(context.Background(), ListUsersParams{
// 		Limit:  5,
// 		Offset: 0,
// 	})

// 	require.NoError(t, err)
// 	require.NotEmpty(t, users)
// 	require.Equal(t, len(users), 5)
// }
