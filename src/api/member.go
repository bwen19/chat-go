package api

import (
	"context"
	"gochat/src/db"
	"gochat/src/hub"
)

// ======================== // addMembers // ======================== //

type AddMembersRequest struct {
	RoomID    int64   `json:"room_id"`
	MemberIDs []int64 `json:"member_ids"`
}
type AddMembersResponse struct {
	RoomID  int64            `json:"room_id"`
	Members []*db.MemberInfo `json:"members"`
}

func (s *Server) addMembers(ctx context.Context, client *hub.Client, req *AddMembersRequest) error {
	return nil
}

// ======================== // deleteMembers // ======================== //

type DeleteMembersRequest struct {
	RoomID    int64   `json:"room_id"`
	MemberIDs []int64 `json:"member_ids"`
}
type DeleteMembersResponse struct {
	RoomID    int64   `json:"room_id"`
	MemberIDs []int64 `json:"member_ids"`
}

func (s *Server) deleteMembers(ctx context.Context, client *hub.Client, req *DeleteMembersRequest) error {
	return nil
}
