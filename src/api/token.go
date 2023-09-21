package api

import (
	"database/sql"
	"errors"
	"gochat/src/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) verifyRefreshToken(c *gin.Context, refreshToken string) (*utils.Payload, error) {
	payload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	session, err := s.store.GetSession(c, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("session not exists")
		}
		return nil, errors.New("failed to get session")
	}
	if session.UserID != payload.UserID {
		return nil, errors.New("mismatched session user")
	}
	if session.RefreshToken != refreshToken {
		return nil, errors.New("mismatched session token")
	}

	return payload, nil
}
