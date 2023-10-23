package db

import (
	"context"
	"gochat/src/db/sqlc"
	"gochat/src/util"
)

type CreateUserParams struct {
	Username string
	Password string
	Role     string
}

func (s *dbStore) CreateUser(ctx context.Context, arg *CreateUserParams) (*UserInfo, error) {
	var userInfo *UserInfo

	hashedPassword, err := util.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}

	err = s.execTx(ctx, func(q *sqlc.Queries) error {
		room, err := q.InsertRoom(ctx, &sqlc.InsertRoomParams{
			Name:     personalRoomName,
			Cover:    personalRoomCover,
			Category: CategoryPersonal,
		})
		if err != nil {
			return err
		}

		user, err := q.InsertUser(ctx, &sqlc.InsertUserParams{
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

		err = q.InsertMember(ctx, &sqlc.InsertMemberParams{
			RoomID:   room.ID,
			MemberID: user.ID,
			Rank:     RankMember,
		})
		if err != nil {
			return err
		}

		userInfo = NewUserInfo(user)
		return nil
	})

	return userInfo, err
}

func (s *dbStore) GetUserByID(ctx context.Context, userID int64) (*UserInfo, error) {
	userInfo := &UserInfo{}

	if err := s.GetCache(ctx, userKey(userID), userInfo); err != nil {
		if err == ErrRedisNil {
			user, err := s.RetrieveUserByID(ctx, userID)
			if err != nil {
				return nil, err
			}

			userInfo = NewUserInfo(user)
			if err = s.cacheUser(ctx, userInfo); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return userInfo, nil
}

func (s *dbStore) GetUserByName(ctx context.Context, username string) (*sqlc.User, error) {
	user, err := s.RetrieveUserByName(ctx, username)
	if err != nil {
		return nil, err
	}

	userInfo := NewUserInfo(user)
	if err = s.cacheUser(ctx, userInfo); err != nil {
		return nil, err
	}

	return user, nil
}

type ListUsersParams sqlc.RetrieveUsersParams

func (s *dbStore) ListUsers(ctx context.Context, arg *ListUsersParams) (int64, []*UserInfo, error) {
	var total int64

	users, err := s.RetrieveUsers(ctx, (*sqlc.RetrieveUsersParams)(arg))
	if err != nil {
		return total, nil, err
	}

	userList := make([]*UserInfo, 0, 5)

	if len(users) > 0 {
		total = users[0].Total

		for _, u := range users {
			userInfo := &UserInfo{
				ID:       u.ID,
				Username: u.Username,
				Nickname: u.Nickname,
				Avatar:   u.Avatar,
				Role:     u.Role,
				Deleted:  u.Deleted,
				CreateAt: u.CreateAt,
			}
			userList = append(userList, userInfo)
		}
	}

	return total, userList, nil
}

type ModifyUserParams sqlc.UpdateUserParams

func (s *dbStore) ModifyUser(ctx context.Context, arg *ModifyUserParams) (*UserInfo, error) {
	user, err := s.UpdateUser(ctx, (*sqlc.UpdateUserParams)(arg))
	if err != nil {
		return nil, err
	}

	userInfo := NewUserInfo(user)
	if err = s.cacheUser(ctx, userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *dbStore) RemoveUser(ctx context.Context, userID int64) error {
	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		err := q.DeleteSessionsByUser(ctx, userID)
		if err != nil {
			return err
		}

		roomIDs, err := q.RetrieveOwnerRoomIDs(ctx, userID)
		if err != nil {
			return err
		}

		friendRoomIDs, err := q.DeleteFriendsByUser(ctx, userID)
		if err != nil {
			return err
		}

		roomIDs = append(roomIDs, friendRoomIDs...)
		err = q.DeleteMessagesByUser(ctx, &sqlc.DeleteMessagesByUserParams{
			UserID:  userID,
			RoomIds: roomIDs,
		})
		if err != nil {
			return err
		}

		err = q.DeleteMembersByUser(ctx, &sqlc.DeleteMembersByUserParams{
			UserID:  userID,
			RoomIds: roomIDs,
		})
		if err != nil {
			return err
		}

		personalRoomID, err := q.DeleteUser(ctx, userID)
		if err != nil {
			return err
		}

		roomIDs = append(roomIDs, personalRoomID)
		if err = q.DeleteRooms(ctx, roomIDs); err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	return s.DelCache(ctx, userKey(userID))
}

func (s *dbStore) cacheUser(ctx context.Context, user *UserInfo) error {
	return s.SetCache(ctx, userKey(user.ID), user)
}
