package main

import (
	"context"
	"gochat/src/api"
	"gochat/src/db"
	"gochat/src/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	store, err := db.NewStore(&config)
	if err != nil {
		log.Fatal("failed to create store")
	}
	log.Println("db connected successfully")

	server, err := api.NewServer(&config, store)
	if err != nil {
		log.Fatal("cannot create api server")
	}

	srv := server.SetupHttpServer()
	log.Printf("start Http server at %s", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("server exiting")
}
