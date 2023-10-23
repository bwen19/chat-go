package api

import (
	"context"
	"gochat/src/db"
	"gochat/src/hub"
	"gochat/src/util"
	"gochat/src/util/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type Server struct {
	srv        *http.Server
	config     *util.Config
	cron       *cron.Cron
	hub        *hub.Hub
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config *util.Config) (*Server, error) {
	hub := hub.NewHub()

	store, err := db.New(config)
	if err != nil {
		return nil, err
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		hub:        hub,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupCronJob()
	server.setupHttpServer()
	return server, nil
}

func (s *Server) Start() error {
	s.cron.Start()
	s.hub.Start()
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.cron.Stop()
	s.store.DumpAllMessages(ctx)
	s.hub.Stop()
	return s.srv.Shutdown(ctx)
}

func (s *Server) setupCronJob() {
	c := cron.New()
	c.AddFunc("0 3 * * *", s.store.DumpPartialMessages(context.Background()))
	s.cron = c
}

func (s *Server) setupHttpServer() {
	router := gin.Default()

	router.GET("/ws/:authorization", s.handleWebsocket)

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
	auth.GET("/user/username", s.findUser)

	admin := auth.Use(s.adminMiddleware())
	admin.POST("/user", s.createUser)
	admin.DELETE("/user/:user_id", s.deleteUser)
	admin.GET("/user", s.listUsers)
	admin.PATCH("/user", s.updateUser)

	s.srv = &http.Server{
		Addr:    s.config.ServerAddress,
		Handler: router,
	}
}
