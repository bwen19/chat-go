package main

import (
	"context"
	"gochat/src/api"
	"gochat/src/core"
	"gochat/src/util"
	"gochat/src/ws"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	state, err := core.NewState(&config)
	if err != nil {
		log.Fatal("failed to initialize:", err)
	}
	defer state.Close()

	// register router
	router := gin.Default()
	apiServer := api.NewServer(state)
	wsServer := ws.NewServer(state)
	apiServer.RegisterRouter(router)
	wsServer.RegisterRouter(router)

	// start HTTP server
	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to serve: ", err)
		}
	}()
	log.Printf("start HTTP server at %s", config.ServerAddress)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}
}
