package api

import (
	"fmt"
	"gochat/src/db"
	"gochat/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      db.Store
	config     *utils.Config
	tokenMaker utils.TokenMaker
}

func NewServer(config *utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := utils.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	return server, nil
}

func (s *Server) SetupHttpServer() *http.Server {
	router := s.newRouter()

	return &http.Server{
		Addr:    s.config.ServerAddress,
		Handler: router,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
