package db

import (
	"context"
	"gochat/src/util"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	oldUser := createRandomUser(t)

	newName := util.RandomString(8)
	newNick := util.RandomNumString(8)
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.User.ID,
		Username: pgtype.Text{
			String: newName,
			Valid:  true,
		},
		Nickname: pgtype.Text{
			String: newNick,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, oldUser.User.ID, updatedUser.ID)
	require.Equal(t, newName, updatedUser.Username)
	require.Equal(t, newNick, updatedUser.Nickname)
	require.Equal(t, oldUser.User.HashedPassword, updatedUser.HashedPassword)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	users, err := testStore.ListUsers(context.Background(), ListUsersParams{
		Limit:  5,
		Offset: 0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, len(users), 5)
}
