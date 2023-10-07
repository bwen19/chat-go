package api

import (
	"gochat/src/util/state"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*state.State
}

func NewServer(state *state.State) *Server {
	return &Server{
		State: state,
	}
}

func (s *Server) RegisterRouter(router *gin.Engine) *gin.Engine {
	api := router.Group("/api")
	api.POST("/login", s.login)
	api.POST("/auto-login", s.autoLogin)
	api.POST("/renew-token", s.renewToken)
	api.POST("/logout", s.logout)

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

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
