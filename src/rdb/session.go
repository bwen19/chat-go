package rdb

import (
	"context"
	"fmt"
	db "gochat/src/db/sqlc"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	UserID       int64  `redis:"user_id"`
	RefreshToken string `redis:"refresh_token"`
	ClientIp     string `redis:"client_ip"`
	UserAgent    string `redis:"user_agent"`
}

type SessionCache interface {
	SetSession(ctx context.Context, ID uuid.UUID, sess *db.Session) error
	GetSession(ctx context.Context, ID uuid.UUID) (db.Session, error)
	DeleteSession(ctx context.Context, ID uuid.UUID) error
}

func sessionKey(ID uuid.UUID) string {
	return fmt.Sprintf("session:%s", ID)
}

func (r *Redis) SetSession(ctx context.Context, ID uuid.UUID, sess *db.Session) error {
	k := sessionKey(ID)
	v := Session{
		UserID:       sess.UserID,
		RefreshToken: sess.RefreshToken,
		ClientIp:     sess.ClientIp,
		UserAgent:    sess.UserAgent,
	}
	if err := r.rdb.HSet(ctx, k, v).Err(); err != nil {
		return err
	}
	if err := r.rdb.Expire(ctx, k, 24*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetSession(ctx context.Context, ID uuid.UUID) (db.Session, error) {
	ret := db.Session{}
	cmd := r.rdb.HGetAll(ctx, sessionKey(ID))
	if err := cmd.Err(); err != nil {
		return ret, err
	}
	if len(cmd.Val()) == 0 {
		return ret, redis.Nil
	}
	sess := &Session{}
	if err := cmd.Scan(sess); err != nil {
		return ret, err
	}

	ret.ID = ID
	ret.UserID = sess.UserID
	ret.RefreshToken = sess.RefreshToken
	ret.ClientIp = sess.ClientIp
	ret.UserAgent = sess.UserAgent

	return ret, nil
}

func (r *Redis) DeleteSession(ctx context.Context, ID uuid.UUID) error {
	err := r.rdb.Del(ctx, sessionKey(ID)).Err()
	if err != nil {
		return err
	}
	return nil
}
