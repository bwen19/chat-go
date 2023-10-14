package core

import (
	"context"
	"encoding/json"
	"fmt"
	db "gochat/src/db/sqlc"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserInfo struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Nickname string    `json:"nickname"`
	Role     string    `json:"role"`
	RoomID   int64     `json:"room_id"`
	Deleted  bool      `json:"deleted"`
	CreateAt time.Time `json:"create_at"`
}

func NewUserInfo(v *db.User) *UserInfo {
	return &UserInfo{
		ID:       v.ID,
		Username: v.Username,
		Nickname: v.Nickname,
		Avatar:   v.Avatar,
		Role:     v.Role,
		Deleted:  v.Deleted,
		RoomID:   v.RoomID,
		CreateAt: v.CreateAt,
	}
}

func (u *UserInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *UserInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func userKey(userID int64) string {
	return fmt.Sprintf("user:%d", userID)
}

func (s *State) GetUser(ctx context.Context, userID int64) (*UserInfo, error) {
	key := userKey(userID)
	res := &UserInfo{}

	if err := s.Cache.Get(ctx, key, res); err != nil {
		if err == redis.Nil {
			user, err := s.Store.GetUser(ctx, userID)
			if err != nil {
				return nil, err
			}

			res = NewUserInfo(&user)
			if err = s.Cache.Set(ctx, key, res); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (s *State) CacheUser(ctx context.Context, user *db.User) error {
	val := NewUserInfo(user)
	return s.Cache.Set(ctx, userKey(val.ID), val)
}
