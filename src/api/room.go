package api

import (
	"context"
	"gochat/src/db"
	"gochat/src/hub"
)

// ======================== // newRoom // ======================== //

type NewRoomRequest struct {
	Name      string  `json:"name"`
	MemberIDs []int64 `json:"member_ids"`
}
type NewRoomResponse struct {
	Room *db.RoomInfo `json:"room"`
}

func (s *Server) newRoom(ctx context.Context, client *hub.Client, req *NewRoomRequest) error {
	return nil
}

// ======================== // updateRoom // ======================== //

type UpdateRoomResquest struct {
	RoomID int64  `json:"room_id"`
	Name   string `json:"name"`
}
type UpdateRoomResponse struct {
	RoomID int64  `json:"room_id"`
	Name   string `json:"name"`
}

func (s *Server) updateRoom(ctx context.Context, client *hub.Client, req *UpdateRoomResponse) error {
	return nil
}

// ======================== // deleteRoom // ======================== //

type DeleteRoomRequest struct {
	RoomID int64 `json:"room_id"`
}
type DeleteRoomResponse struct {
	RoomID int64 `json:"room_id"`
}

func (s *Server) deleteRoom(ctx context.Context, client *hub.Client, req *DeleteRoomRequest) error {
	return nil
}

// ======================== // leaveRoom // ======================== //

type LeaveRoomRequest struct {
	RoomID int64 `json:"room_id"`
}

func (s *Server) leaveRoom(ctx context.Context, client *hub.Client, req *LeaveRoomRequest) error {
	return nil
}
