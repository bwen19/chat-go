package db

import (
	"context"
	"gochat/src/util"
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

func (s *SqlStore) CreateUserTx(c context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(c, func(q *Queries) error {
		hashedPassword, err := util.HashPassword(arg.Password)
		if err != nil {
			return err
		}

		room, err := q.CreateRoom(c, CreateRoomParams{
			Name:     defaultRoomName,
			Cover:    defaultCover,
			Category: personalCategory,
		})
		if err != nil {
			return err
		}

		result.User, err = q.CreateUser(c, CreateUserParams{
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

		err = q.CreateRoomMember(c, CreateRoomMemberParams{
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

func (s *SqlStore) DeleteUserTx(c context.Context, userID int64) error {
	err := s.execTx(c, func(q *Queries) error {
		if err := q.DeleteSessionByUser(c, userID); err != nil {
			return err
		}

		roomIds, err := q.DeleteFriendByUser(c, userID)
		if err != nil {
			return err
		}

		err = q.DeleteMessageByUser(c, DeleteMessageByUserParams{
			UserID:  userID,
			RoomIds: roomIds,
		})
		if err != nil {
			return err
		}

		err = q.DeleteMemberByUser(c, DeleteMemberByUserParams{
			UserID:  userID,
			RoomIds: roomIds,
		})
		if err != nil {
			return err
		}

		r2, err := q.DeleteUser(c, userID)
		if err != nil {
			return err
		}

		roomIds = append(roomIds, r2...)
		if err = q.DeleteRooms(c, roomIds); err != nil {
			return err
		}

		return err
	})

	return err
}
