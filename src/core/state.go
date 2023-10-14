package core

import (
	"fmt"
	db "gochat/src/db/sqlc"
	"gochat/src/rdb"
	"gochat/src/util"
	"gochat/src/util/token"
	"log"
)

type State struct {
	Config     *util.Config
	Store      db.Store
	Cache      rdb.Cache
	TokenMaker token.Maker
}

func NewState(config *util.Config) (*State, error) {
	store, err := db.NewStore(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}
	log.Print("connect to Postgres")

	cache, err := rdb.NewRedis(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}
	log.Print("connect to Redis")

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	state := &State{
		Config:     config,
		Store:      store,
		Cache:      cache,
		TokenMaker: tokenMaker,
	}

	state.runCron()
	return state, nil
}

func (s *State) Close() {
	s.saveAllMessages()
	s.Store.Close()
	s.Cache.Close()
}
