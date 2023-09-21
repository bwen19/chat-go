package api

import (
	"errors"
	"gochat/src/db"
	"gochat/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ======================== // Login // ======================== //

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,alphanum"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUserByName(c, req.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = utils.CheckPassword(req.Username, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Role,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Role,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.store.CreateSession(c, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ClientIp:     c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		ExpireAt:     refreshPayload.ExpireAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := LoginResponse{
		User:         convertUser(&user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, rsp)
}

// ======================== // AutoLogin // ======================== //

type AutoLoginRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AutoLoginResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

func (s *Server) AutoLogin(c *gin.Context) {
	var req AutoLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.verifyRefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(c, payload.UserID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.Deleted {
		err = errors.New("the user is not available")
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	accessToken, _, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Role,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := &AutoLoginResponse{
		User:        convertUser(&user),
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, rsp)
}

// ======================== // RenewToken // ======================== //

func (s *Server) RenewToken(*gin.Context) {

}

// ======================== // Logout // ======================== //

func (s *Server) Logout(*gin.Context) {

}
