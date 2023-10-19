package db

import (
	"context"
	"gochat/src/db/sqlc"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ClientIp     string    `json:"client_ip"`
	UserAgent    string    `json:"user_agent"`
	ExpireAt     time.Time `json:"expire_at"`
}

func (s *dbStore) CreateSession(ctx context.Context, arg *CreateSessionParams) error {
	sess, err := s.InsertSession(ctx, &sqlc.InsertSessionParams{
		ID:           arg.ID,
		UserID:       arg.UserID,
		RefreshToken: arg.RefreshToken,
		ClientIp:     arg.ClientIp,
		UserAgent:    arg.UserAgent,
		ExpireAt:     arg.ExpireAt,
	})
	if err != nil {
		return err
	}

	sessInfo := NewSessionInfo(sess)
	return s.cacheSession(ctx, sessInfo)
}

func (s *dbStore) cacheSession(ctx context.Context, sessInfo *SessionInfo) error {
	return s.SetCache(ctx, sessionKey(sessInfo.ID), sessInfo)
}

func (s *dbStore) GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionInfo, error) {
	sessInfo := &SessionInfo{}

	if err := s.GetCache(ctx, sessionKey(sessionID), sessInfo); err != nil {
		if err == redis.Nil {
			sess, err := s.RetrieveSession(ctx, sessionID)
			if err != nil {
				return nil, err
			}

			sessInfo = NewSessionInfo(sess)
			if s.cacheSession(ctx, sessInfo); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return sessInfo, nil
}

func (s *dbStore) RemoveSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.DeleteSession(ctx, sessionID); err != nil {
		return err
	}
	return s.DelCache(ctx, sessionKey(sessionID))
}
