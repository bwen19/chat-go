package api

import (
	"errors"
	"gochat/src/db"
	"gochat/src/util"
	"gochat/src/util/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ======================== // login // ======================== //

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	IsAdmin  *bool  `json:"is_admin" binding:"required"`
}
type LoginResponse struct {
	User         *db.UserInfo `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func (s *Server) login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	user, err := s.store.GetUserByName(ctx, req.Username)
	if err != nil {
		NotFoundResponse(ctx, err)
		return
	}

	if user.Deleted || (*req.IsAdmin && user.Role != db.RoleAdmin) {
		PermissionDeniedResponse(ctx)
		return
	}

	if err = util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.ID,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}

	err = s.store.CreateSession(ctx, &db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpireAt:     refreshPayload.ExpireAt,
	})
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}

	rsp := &LoginResponse{
		User:         db.NewUserInfo(user),
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
	User        *db.UserInfo `json:"user"`
	AccessToken string       `json:"access_token"`
}

func (s *Server) autoLogin(ctx *gin.Context) {
	var req AutoLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	user, err := s.store.GetUserByID(ctx, payload.UserID)
	if err != nil {
		NotFoundResponse(ctx, err)
		return
	}

	if user.Deleted {
		PermissionDeniedResponse(ctx)
		return
	}

	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}

	rsp := &AutoLoginResponse{
		User:        user,
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
		InvalidArgumentResponse(ctx)
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	accessToken, _, err := s.tokenMaker.CreateToken(
		payload.UserID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		InternalErrorResponse(ctx)
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
		InvalidArgumentResponse(ctx)
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	if err = s.store.RemoveSession(ctx, payload.ID); err != nil {
		InternalErrorResponse(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

// ======================== // verifyRefreshToken // ======================== //

func (s *Server) verifyRefreshToken(ctx *gin.Context, refreshToken string) (*token.Payload, error) {
	payload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	session, err := s.store.GetSession(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if session.UserID != payload.UserID || session.RefreshToken != refreshToken {
		return nil, errors.New("session does not match")
	}

	return payload, nil
}
