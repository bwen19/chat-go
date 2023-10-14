package api

import (
	"errors"
	"gochat/src/core"
	db "gochat/src/db/sqlc"
	"gochat/src/util"
	"gochat/src/util/token"
	"log"
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
	User         *core.UserInfo `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
}

func (s *Server) login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.InvalidArgumentResponse(ctx)
		return
	}

	user, err := s.Store.GetUserByName(ctx, req.Username)
	if err != nil {
		util.RecordNotFoundResponse(ctx, err)
		return
	}

	if err = s.CacheUser(ctx, &user); err != nil {
		util.InternalErrorResponse(ctx)
		log.Println("err:", err)
		return
	}

	if user.Deleted || (*req.IsAdmin && user.Role != "admin") {
		util.PermissionDeniedResponse(ctx)
		return
	}

	if err = util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	accessToken, _, err := s.TokenMaker.CreateToken(
		user.ID,
		s.Config.AccessTokenDuration,
	)
	if err != nil {
		util.InternalErrorResponse(ctx)
		return
	}

	refreshToken, refreshPayload, err := s.TokenMaker.CreateToken(
		user.ID,
		s.Config.RefreshTokenDuration,
	)
	if err != nil {
		util.InternalErrorResponse(ctx)
		return
	}

	sess, err := s.Store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpireAt:     refreshPayload.ExpireAt,
	})
	if err != nil {
		util.InternalErrorResponse(ctx)
		return
	}

	if err = s.CacheSession(ctx, &sess); err != nil {
		util.InternalErrorResponse(ctx)
		log.Println("err:", err)
		return
	}

	rsp := &LoginResponse{
		User:         core.NewUserInfo(&user),
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
	User        *core.UserInfo `json:"user"`
	AccessToken string         `json:"access_token"`
}

func (s *Server) autoLogin(ctx *gin.Context) {
	var req AutoLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.InvalidArgumentResponse(ctx)
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	user, err := s.GetUser(ctx, payload.UserID)
	if err != nil {
		util.RecordNotFoundResponse(ctx, err)
		return
	}

	if user.Deleted {
		util.PermissionDeniedResponse(ctx)
		return
	}

	accessToken, _, err := s.TokenMaker.CreateToken(
		user.ID,
		s.Config.AccessTokenDuration,
	)
	if err != nil {
		util.InternalErrorResponse(ctx)
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
		util.InvalidArgumentResponse(ctx)
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	accessToken, _, err := s.TokenMaker.CreateToken(
		payload.UserID,
		s.Config.AccessTokenDuration,
	)
	if err != nil {
		util.InternalErrorResponse(ctx)
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
		util.InvalidArgumentResponse(ctx)
		return
	}

	payload, err := s.verifyRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	err = s.Store.DeleteSession(ctx, db.DeleteSessionParams{
		ID:     payload.ID,
		UserID: payload.UserID,
	})
	if err != nil {
		util.InternalErrorResponse(ctx)
		return
	}

	if err = s.DelSession(ctx, payload.ID); err != nil {
		util.InternalErrorResponse(ctx)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// ======================== // utils // ======================== //

func (s *Server) verifyRefreshToken(ctx *gin.Context, refreshToken string) (*token.Payload, error) {
	payload, err := s.TokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	session, err := s.GetSession(ctx, payload.ID)
	if err != nil {
		return nil, err
	}

	if session.UserID != payload.UserID || session.RefreshToken != refreshToken {
		return nil, errors.New("session does not match")
	}

	return payload, nil
}
