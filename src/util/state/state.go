package state

import (
	"fmt"
	db "gochat/src/db/sqlc"
	"gochat/src/rdb"
	"gochat/src/util"
	"gochat/src/util/token"
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

	cache, err := rdb.NewRedis(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	state := &State{
		Config:     config,
		Store:      store,
		Cache:      cache,
		TokenMaker: tokenMaker,
	}
	return state, nil
}
