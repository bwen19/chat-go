package api

import (
	"errors"
	db "gochat/src/db/sqlc"
	"gochat/src/util"
	"gochat/src/util/token"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// ======================== // login // ======================== //

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}
type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.Store.GetUserByName(ctx, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = s.Cache.SetUser(ctx, user.ID, &user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Deleted {
		err = errors.New("the user is not available")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	if err = util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := s.TokenMaker.CreateToken(
		user.ID,
		s.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.TokenMaker.CreateToken(
		user.ID,
		s.Config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpireAt:     refreshPayload.ExpireAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := &LoginResponse{
		User:         convertUser(&user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}

// ======================== // autoLogin // ======================== //

type AutoLoginRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type AutoLoginResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
}

func (s *Server) autoLogin(ctx *gin.Context) {
	var req AutoLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := s.GetUser(ctx, payload.UserID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Deleted {
		err = errors.New("the user is not available")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	accessToken, _, err := s.TokenMaker.CreateToken(
		user.ID,
		s.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := &AutoLoginResponse{
		User:        convertUser(&user),
		AccessToken: accessToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}

// ======================== // renewToken // ======================== //

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type RenewTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *Server) renewToken(ctx *gin.Context) {
	var req RenewTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := s.TokenMaker.CreateToken(
		payload.UserID,
		s.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := &RenewTokenResponse{
		AccessToken: accessToken,
	}
	ctx.JSON(http.StatusOK, rsp)
}

// ======================== // logout // ======================== //

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Server) logout(ctx *gin.Context) {
	var req LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = s.Store.DeleteSession(ctx, db.DeleteSessionParams{
		ID:     payload.ID,
		UserID: payload.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.Cache.DeleteSession(ctx, payload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// ======================== // utils // ======================== //

func (s *Server) verifyRefreshToken(ctx *gin.Context, refreshToken string) (*token.Payload, error) {
	payload, err := s.TokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	session, err := s.Cache.GetSession(ctx, payload.ID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			session, err = s.Store.GetSession(ctx, payload.ID)
			log.Println("get session from db")
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					return nil, errors.New("session not exists")
				}
				return nil, errors.New("failed to get session")
			}

			if err = s.Cache.SetSession(ctx, payload.ID, &session); err != nil {
				return nil, errors.New("failed to store session in cache")
			}
		} else {
			return nil, errors.New("internal redis error")
		}
	}

	if session.UserID != payload.UserID {
		return nil, errors.New("mismatched session user")
	}

	if session.RefreshToken != refreshToken {
		return nil, errors.New("mismatched session token")
	}

	return payload, nil
}
