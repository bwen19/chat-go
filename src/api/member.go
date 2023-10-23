package api

import (
	"context"
	"errors"
	"gochat/src/db"
	"gochat/src/db/sqlc"
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
	userID := client.GetUserID()
	roomID := req.RoomID

	member, err := s.store.RetrieveMember(ctx, &sqlc.RetrieveMemberParams{
		RoomID:   roomID,
		MemberID: userID,
	})
	if err != nil {
		return errInternal
	}

	if member.Rank == db.RankMember {
		return errDenied
	}

	members, err := s.store.AddMembers(ctx, roomID, req.MemberIDs)
	if err != nil {
		return errors.New("error: add members")
	}

	rsp0 := &AddMembersResponse{RoomID: roomID, Members: members}
	evt := newWebsocketEvent("add-members", rsp0)
	s.hub.BroadcastToRoom(evt, roomID)

	room, err := s.store.GetRoomInfo(ctx, roomID)
	if err != nil {
		return errInternal
	}

	userIDs := make([]int64, 0, len(members))
	for _, m := range members {
		userIDs = append(userIDs, m.ID)
	}
	s.hub.JoinRoom(roomID, userIDs...)

	rsp1 := &NewRoomResponse{Room: room}
	evt = newWebsocketEvent("new-room", rsp1)
	s.hub.BroadcastToUsers(evt, userIDs...)

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
	userID := client.GetUserID()
	roomID := req.RoomID

	member, err := s.store.RetrieveMember(ctx, &sqlc.RetrieveMemberParams{
		RoomID:   roomID,
		MemberID: userID,
	})
	if err != nil {
		return errInternal
	}

	if member.Rank == db.RankMember {
		return errDenied
	}

	memberIDs, err := s.store.DeleteMembers(ctx, &sqlc.DeleteMembersParams{
		RoomID:    roomID,
		MemberIds: req.MemberIDs,
	})
	if err != nil {
		return errInternal
	}

	rsp := &DeleteMembersResponse{RoomID: roomID, MemberIDs: memberIDs}
	evt := newWebsocketEvent("delete-members", rsp)
	s.hub.BroadcastToRoom(evt, roomID)

	return nil
}
