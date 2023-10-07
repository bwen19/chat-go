package rdb

import (
	"context"
	"fmt"
	db "gochat/src/db/sqlc"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	Username string `redis:"username"`
	Avatar   string `redis:"avatar"`
	Nickname string `redis:"nickname"`
	Role     string `redis:"role"`
	RoomID   int64  `redis:"room_id"`
	Deleted  bool   `redis:"deleted"`
}

type UserCache interface {
	SetUser(ctx context.Context, userID int64, user *db.User) error
	GetUser(ctx context.Context, userID int64) (db.User, error)
	DeleteUser(ctx context.Context, userID int64) error
}

func userKey(userID int64) string {
	return fmt.Sprintf("user:%d", userID)
}

func (r *Redis) SetUser(ctx context.Context, userID int64, user *db.User) error {
	k := userKey(userID)
	v := User{
		Username: user.Username,
		Avatar:   user.Avatar,
		Nickname: user.Nickname,
		Role:     user.Role,
		RoomID:   user.RoomID,
		Deleted:  user.Deleted,
	}
	if err := r.rdb.HSet(ctx, k, v).Err(); err != nil {
		return err
	}
	if err := r.rdb.Expire(ctx, k, 24*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetUser(ctx context.Context, userID int64) (db.User, error) {
	ret := db.User{}
	cmd := r.rdb.HGetAll(ctx, userKey(userID))
	if err := cmd.Err(); err != nil {
		return ret, err
	}
	if len(cmd.Val()) == 0 {
		return ret, redis.Nil
	}
	user := &User{}
	if err := cmd.Scan(user); err != nil {
		return ret, err
	}

	ret.Username = user.Username
	ret.Avatar = user.Avatar
	ret.Nickname = user.Nickname
	ret.Role = user.Role
	ret.RoomID = user.RoomID
	ret.Deleted = user.Deleted

	return ret, nil
}

func (r *Redis) DeleteUser(ctx context.Context, userID int64) error {
	err := r.rdb.Del(ctx, userKey(userID)).Err()
	if err != nil {
		return err
	}
	return nil
}
