package api

import (
	"gochat/src/core"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*core.State
}

func NewServer(state *core.State) *Server {
	return &Server{
		State: state,
	}
}

func (s *Server) RegisterRouter(router *gin.Engine) *gin.Engine {
	api := router.Group("/api")
	api.POST("/auth/login", s.login)
	api.POST("/auth/auto-login", s.autoLogin)
	api.POST("/auth/renew-token", s.renewToken)
	api.POST("/auth/logout", s.logout)

	auth := api.Use(s.authMiddleware())
	auth.PATCH("/user/password", s.changePassword)
	auth.PATCH("/user/info", s.changeUserInfo)
	auth.DELETE("/session/:session_id", s.deleteSession)
	auth.GET("/session", s.listSessions)

	admin := auth.Use(s.adminMiddleware())
	admin.POST("/user", s.createUser)
	admin.DELETE("/user/:user_id", s.deleteUser)
	admin.GET("/user", s.listUsers)
	admin.PATCH("/user", s.updateUser)

	return router
}
