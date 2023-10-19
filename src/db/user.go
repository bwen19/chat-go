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

		err = q.InsertRoomMember(ctx, &sqlc.InsertRoomMemberParams{
			RoomID:   room.ID,
			MemberID: user.ID,
			Rank:     RankMember,
		})
		if err != nil {
			return err
		}

		userInfo = NewUserInfo(user)
		return err
	})

	return userInfo, err
}

func (s *dbStore) cacheUser(ctx context.Context, user *UserInfo) error {
	return s.SetCache(ctx, userKey(user.ID), user)
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
		return user, err
	}

	userInfo := NewUserInfo(user)
	if err = s.cacheUser(ctx, userInfo); err != nil {
		return user, err
	}

	return user, nil
}

func (s *dbStore) GetUsers(ctx context.Context, arg *sqlc.ListUsersParams) (int64, []*UserInfo, error) {
	var total int64

	users, err := s.ListUsers(ctx, arg)
	if err != nil {
		return total, nil, err
	}

	userList := make([]*UserInfo, 0, 5)

	if len(users) > 0 {
		total = users[0].Total

		for _, user := range users {
			userList = append(userList, &UserInfo{
				ID:       user.ID,
				Username: user.Username,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
				Role:     user.Role,
				Deleted:  user.Deleted,
				CreateAt: user.CreateAt,
			})
		}
	}

	return total, userList, nil
}

func (s *dbStore) ModifyUser(ctx context.Context, arg *sqlc.UpdateUserParams) (*UserInfo, error) {
	user, err := s.UpdateUser(ctx, arg)
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
		err := q.DeleteSessionByUser(ctx, userID)
		if err != nil {
			return err
		}

		roomIds, err := q.DeleteFriendByUser(ctx, userID)
		if err != nil {
			return err
		}

		err = q.DeleteMessageByUser(ctx, &sqlc.DeleteMessageByUserParams{
			UserID:  userID,
			RoomIds: roomIds,
		})
		if err != nil {
			return err
		}

		err = q.DeleteMemberByUser(ctx, &sqlc.DeleteMemberByUserParams{
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
	if err != nil {
		return err
	}

	return s.DelCache(ctx, userKey(userID))
}
