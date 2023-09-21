package api

import (
	"gochat/src/db"
)

func convertUser(user *db.User) User {
	return User{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     user.Role,
		Deleted:  user.Deleted,
		RoomID:   user.RoomID,
		CreateAt: user.CreateAt,
	}
}
