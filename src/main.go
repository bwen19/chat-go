package main

import (
	"context"
	"gochat/src/api"
	"gochat/src/util"
	"gochat/src/util/state"
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
		log.Fatal("failed to load config: ", err)
	}

	srv, err := newHttpServer(&config)
	if err != nil {
		log.Fatal("failed to create HTTP server: ", err)
	}

	log.Printf("start HTTP server at %s", config.ServerAddress)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	<-ctx.Done()
	log.Println("server exiting")
}

func newHttpServer(config *util.Config) (*http.Server, error) {
	state, err := state.NewState(config)
	if err != nil {
		return nil, err
	}
	log.Println("initialization complete")

	server := api.NewServer(state)
	wsServer := ws.NewServer(state)

	router := gin.Default()
	server.RegisterRouter(router)
	wsServer.RegisterRouter(router)

	httpServer := &http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}
	return httpServer, nil
}
