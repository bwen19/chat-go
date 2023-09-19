package api

import (
	"fmt"
	"gochat/src/db"
	"gochat/src/pb"
	"gochat/src/utils"
)

type Server struct {
	pb.UnimplementedChatServer
	config     utils.Config
	store      db.Store
	tokenMaker utils.TokenMaker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := utils.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
