package db

import (
	"context"
	"gochat/src/utils"
)

const (
	defaultRoomName  = "My Room"
	defaultAvatar    = "/cover/personal"
	defaultCover     = "/cover/public"
	defaultRank      = "owner"
	personalCategory = "personal"
)

type CreateUserTxParams struct {
	Username string
	Password string
	Role     string
}

type CreateUserTxResult struct {
	User User
}

func (store *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		hashedPassword, err := utils.HashPassword(arg.Password)
		if err != nil {
			return err
		}

		room, err := q.CreateRoom(ctx, CreateRoomParams{
			Name:     defaultRoomName,
			Cover:    defaultCover,
			Category: personalCategory,
		})
		if err != nil {
			return err
		}

		result.User, err = q.CreateUser(ctx, CreateUserParams{
			Username:       arg.Username,
			HashedPassword: hashedPassword,
			Nickname:       arg.Username,
			Avatar:         defaultAvatar,
			Role:           arg.Role,
			RoomID:         room.ID,
		})
		if err != nil {
			return err
		}

		err = q.CreateRoomMember(ctx, CreateRoomMemberParams{
			RoomID:   room.ID,
			MemberID: result.User.ID,
			Rank:     defaultRank,
		})
		if err != nil {
			return err
		}

		return err
	})

	return result, err
}

func (store *SqlStore) DeleteUserTx(ctx context.Context, userID int64) error {
	err := store.execTx(ctx, func(q *Queries) error {
		if err := q.DeleteSessionByUser(ctx, userID); err != nil {
			return err
		}

		roomIds, err := q.DeleteFriendByUser(ctx, userID)
		if err != nil {
			return err
		}

		err = q.DeleteMessageByUser(ctx, DeleteMessageByUserParams{
			UserID:  userID,
			RoomIds: roomIds,
		})
		if err != nil {
			return err
		}

		err = q.DeleteMemberByUser(ctx, DeleteMemberByUserParams{
			UserID:  userID,
			RoomIds: roomIds,
		})
		if err != nil {
			return err
		}

		r2, err := q.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}

		roomIds = append(roomIds, r2...)
		if err = q.DeleteRooms(ctx, roomIds); err != nil {
			return err
		}

		return err
	})

	return err
}
