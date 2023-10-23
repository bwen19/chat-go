package api

import (
	"context"
	"errors"
	"gochat/src/db"
	"gochat/src/db/sqlc"
	"gochat/src/hub"
	"gochat/src/util"

	"github.com/jackc/pgx/v5/pgtype"
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
	userID := client.GetUserID()
	memberIDs := req.MemberIDs

	room, err := s.store.NewRoom(ctx, req.Name, userID, memberIDs)
	if err != nil {
		return errors.New("error new room")
	}

	memberIDs = append(memberIDs, userID)
	s.hub.JoinRoom(room.ID, memberIDs...)

	rsp := &NewRoomResponse{Room: room}
	evt := newWebsocketEvent("new-room", rsp)
	s.hub.BroadcastToRoom(evt, room.ID)

	return nil
}

// ======================== // updateRoom // ======================== //

type UpdateRoomRequest struct {
	RoomID int64   `json:"room_id"`
	Name   *string `json:"name"`
	Cover  *string `json:"cover"`
}
type UpdateRoomResponse struct {
	RoomID int64  `json:"room_id"`
	Name   string `json:"name"`
	Cover  string `json:"cover"`
}

func (s *Server) updateRoom(ctx context.Context, client *hub.Client, req *UpdateRoomRequest) error {
	arg := &sqlc.UpdateRoomParams{ID: req.RoomID}
	if req.Name != nil {
		name := *req.Name
		if err := util.ValidateName(name); err != nil {
			return errArgument
		}
		arg.Name = pgtype.Text{String: name, Valid: true}
	}
	if req.Cover != nil {
		arg.Cover = pgtype.Text{String: *req.Cover, Valid: true}
	}

	member, err := s.store.RetrieveMember(ctx, &sqlc.RetrieveMemberParams{
		RoomID:   arg.ID,
		MemberID: client.GetUserID(),
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return errNotFound
		}
		return errInternal
	}

	if member.Rank != db.RankOwner && member.Rank != db.RankManager {
		return errDenied
	}

	room, err := s.store.UpdateRoom(ctx, arg)
	if err != nil {
		return errors.New("error: update room")
	}

	rsp := &UpdateRoomResponse{
		RoomID: room.ID,
		Name:   room.Name,
		Cover:  room.Cover,
	}
	evt := newWebsocketEvent("update-room", rsp)
	s.hub.BroadcastToRoom(evt, room.ID)

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
	roomID := req.RoomID

	member, err := s.store.RetrieveMember(ctx, &sqlc.RetrieveMemberParams{
		RoomID:   roomID,
		MemberID: client.GetUserID(),
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return errNotFound
		}
		return errInternal
	}

	if member.Rank != db.RankOwner {
		return errDenied
	}

	err = s.store.RemoveRoom(ctx, roomID)
	if err != nil {
		return errors.New("error: delete room")
	}

	rsp := &DeleteRoomResponse{RoomID: roomID}
	evt := newWebsocketEvent("delete-room", rsp)
	s.hub.BroadcastToRoom(evt, roomID)

	s.hub.DeleteRoom(roomID)

	return nil
}

// ======================== // leaveRoom // ======================== //

type LeaveRoomRequest struct {
	RoomID int64 `json:"room_id"`
}

func (s *Server) leaveRoom(ctx context.Context, client *hub.Client, req *LeaveRoomRequest) error {
	userID := client.GetUserID()
	roomID := req.RoomID

	member, err := s.store.RetrieveMember(ctx, &sqlc.RetrieveMemberParams{
		RoomID:   roomID,
		MemberID: userID,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return errNotFound
		}
		return errInternal
	}

	if member.Rank == db.RankOwner {
		return errDenied
	}

	err = s.store.DeleteMember(ctx, &sqlc.DeleteMemberParams{
		RoomID:   roomID,
		MemberID: userID,
	})
	if err != nil {
		return errors.New("error: leave room")
	}

	s.hub.LeaveRoom(roomID, userID)

	rsp := &DeleteRoomResponse{RoomID: roomID}
	evt := newWebsocketEvent("delete-room", rsp)
	client.SendToSelf(evt)

	rsp1 := &DeleteMembersResponse{
		RoomID:    roomID,
		MemberIDs: []int64{userID},
	}
	evt = newWebsocketEvent("delete-members", rsp1)
	s.hub.BroadcastToRoom(evt, roomID)

	return nil
}
