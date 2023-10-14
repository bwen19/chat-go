package api

import (
	"errors"
	"gochat/src/util/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationKey        = "authorization"
	lowerBearerKey          = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationKey)

		if len(authorizationHeader) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != lowerBearerKey {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		payload, err := s.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			if errors.Is(err, token.ErrExpiredToken) {
				ctx.AbortWithStatus(http.StatusPaymentRequired)
				return
			}
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

		user, err := s.GetUser(ctx, payload.UserID)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if user.Role != "admin" {
			ctx.AbortWithStatus(http.StatusForbidden)
		}
		ctx.Next()
	}
}
