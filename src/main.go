package main

import (
	"context"
	"gochat/src/api"
	"gochat/src/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	server, err := api.NewServer(&config)
	if err != nil {
		log.Fatal("failed to create server: ", err)
	}

	// start HTTP server
	log.Printf("start HTTP server at %s", config.ServerAddress)
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to serve: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}
}
