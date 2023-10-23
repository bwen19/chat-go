package db

import (
	"context"
	"gochat/src/db/sqlc"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CreateSessionParams sqlc.InsertSessionParams

func (s *dbStore) CreateSession(ctx context.Context, arg *CreateSessionParams) error {
	sess, err := s.InsertSession(ctx, (*sqlc.InsertSessionParams)(arg))
	if err != nil {
		return err
	}

	sessInfo := NewSessionInfo(sess)
	return s.cacheSession(ctx, sessInfo)
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

type ListSessionsParams sqlc.RetrieveSessionsParams

func (s *dbStore) ListSessions(ctx context.Context, arg *ListSessionsParams) (int64, []*SessionInfo, error) {
	var total int64

	sessions, err := s.RetrieveSessions(ctx, (*sqlc.RetrieveSessionsParams)(arg))
	if err != nil {
		return total, nil, err
	}

	sessList := make([]*SessionInfo, 0, 5)

	if len(sessions) > 0 {
		total = sessions[0].Total

		for _, session := range sessions {
			sess := &SessionInfo{
				ID:        session.ID,
				ClientIp:  session.ClientIp,
				UserAgent: session.UserAgent,
				ExpireAt:  session.ExpireAt,
				CreateAt:  session.CreateAt,
			}
			sessList = append(sessList, sess)
		}
	}

	return total, sessList, nil
}

func (s *dbStore) RemoveSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.DeleteSession(ctx, sessionID); err != nil {
		return err
	}
	return s.DelCache(ctx, sessionKey(sessionID))
}

func (s *dbStore) cacheSession(ctx context.Context, sessInfo *SessionInfo) error {
	return s.SetCache(ctx, sessionKey(sessInfo.ID), sessInfo)
}
