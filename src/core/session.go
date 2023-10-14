package core

import (
	"context"
	"encoding/json"
	"fmt"
	db "gochat/src/db/sqlc"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionInfo struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ClientIp     string    `json:"client_ip"`
	UserAgent    string    `json:"user_agent"`
}

func NewSessionInfo(v *db.Session) *SessionInfo {
	return &SessionInfo{
		ID:           v.ID,
		UserID:       v.UserID,
		RefreshToken: v.RefreshToken,
		ClientIp:     v.ClientIp,
		UserAgent:    v.UserAgent,
	}
}

func (s *SessionInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *SessionInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func sessionKey(ID uuid.UUID) string {
	return fmt.Sprintf("session:%s", ID)
}

func (s *State) GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionInfo, error) {
	key := sessionKey(sessionID)
	res := &SessionInfo{}

	if err := s.Cache.Get(ctx, key, res); err != nil {
		if err == redis.Nil {
			session, err := s.Store.GetSession(ctx, sessionID)
			if err != nil {
				return nil, err
			}

			res = NewSessionInfo(&session)
			if err = s.Cache.Set(ctx, key, res); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (s *State) CacheSession(ctx context.Context, sess *db.Session) error {
	val := NewSessionInfo(sess)
	return s.Cache.Set(ctx, sessionKey(val.ID), val)
}

func (s *State) DelSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.Cache.Del(ctx, sessionKey(sessionID))
}
