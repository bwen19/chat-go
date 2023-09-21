package api

import "github.com/gin-gonic/gin"

func (s *Server) newRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	api.POST("/login", s.Login)

	auth := api.Use(authMiddleware(s.tokenMaker))
	auth.POST("/user", s.CreateUser)

	return router
}
