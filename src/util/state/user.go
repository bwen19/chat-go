package state

import (
	"context"
	"errors"
	db "gochat/src/db/sqlc"

	"github.com/redis/go-redis/v9"
)

func (s *State) GetUser(ctx context.Context, userID int64) (db.User, error) {
	user, err := s.Cache.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			user, err = s.Store.GetUser(ctx, userID)
			if err != nil {
				return user, err
			}
			if err = s.Cache.SetUser(ctx, user.ID, &user); err != nil {
				return user, err
			}
		} else {
			return user, err
		}
	}
	return user, nil
}
