package api

import (
	"gochat/src/db"
	"gochat/src/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Deleted:   user.Deleted,
		RoomId:    user.RoomID,
		CreatedAt: timestamppb.New(user.CreateAt),
	}
}
