package api

import (
	"context"
	"errors"
	"gochat/src/db"
	"gochat/src/hub"
)

// ======================== // addFriend // ======================== //

type AddFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type AddFriendResponse struct {
	Friend *db.FriendInfo `json:"friend"`
}

func (s *Server) addFriend(ctx context.Context, client *hub.Client, req *AddFriendRequest) error {
	uFriend, fFriend, err := s.store.AddFriend(ctx, client.GetUserID(), req.FriendID)
	if err != nil {
		return errors.New("error adding friend")
	}

	rsp := &AddFriendResponse{Friend: uFriend}
	evt := newWebsocketEvent("add-friend", rsp)
	client.SendToSelf(evt)

	rsp = &AddFriendResponse{Friend: fFriend}
	evt = newWebsocketEvent("add-friend", rsp)
	s.hub.BroadcastToUsers(evt, req.FriendID)

	return nil
}

// ======================== // acceptFriend // ======================== //

type AcceptFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type AcceptFriendResponse struct {
	Friend *db.FriendInfo `json:"friend"`
	Room   *db.RoomInfo   `json:"room"`
}

func (s *Server) acceptFriend(ctx context.Context, client *hub.Client, req *AcceptFriendRequest) error {
	uFriend, uRoom, fFriend, fRoom, err := s.store.AcceptFriend(ctx, client.GetUserID(), req.FriendID)
	if err != nil {
		return errors.New("error accept friend")
	}

	s.hub.JoinRoom(uRoom.ID, uFriend.ID, fFriend.ID)

	rsp := &AcceptFriendResponse{Friend: uFriend, Room: uRoom}
	evt := newWebsocketEvent("accept-friend", rsp)
	client.SendToSelf(evt)

	rsp = &AcceptFriendResponse{Friend: fFriend, Room: fRoom}
	evt = newWebsocketEvent("accept-friend", rsp)
	s.hub.BroadcastToUsers(evt, req.FriendID)

	return nil
}

// ======================== // refuseFriend // ======================== //

type RefuseFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type RefuseFriendResponse struct {
	FriendID int64 `json:"friend_id"`
}

func (s *Server) refuseFriend(ctx context.Context, client *hub.Client, req *RefuseFriendRequest) error {
	userID := client.GetUserID()
	friendID := req.FriendID

	err := s.store.RefuseFriend(ctx, userID, friendID)
	if err != nil {
		return errors.New("error refuse friend")
	}

	rsp := &RefuseFriendResponse{FriendID: friendID}
	evt := newWebsocketEvent("refuse-friend", rsp)
	client.SendToSelf(evt)

	rsp = &RefuseFriendResponse{FriendID: userID}
	evt = newWebsocketEvent("refuse-friend", rsp)
	s.hub.BroadcastToUsers(evt, friendID)

	return nil
}

// ======================== // deleteFriend // ======================== //

type DeleteFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type DeleteFriendResponse struct {
	FriendID int64 `json:"friend_id"`
	RoomID   int64 `json:"room_id"`
}

func (s *Server) deleteFriend(ctx context.Context, client *hub.Client, req *DeleteFriendRequest) error {
	userID := client.GetUserID()
	friendID := req.FriendID

	roomID, err := s.store.RemoveFriend(ctx, userID, friendID)
	if err != nil {
		return errors.New("error delete friend")
	}

	s.hub.LeaveRoom(roomID, userID, req.FriendID)

	rsp := &DeleteFriendResponse{FriendID: friendID, RoomID: roomID}
	evt := newWebsocketEvent("delete-friend", rsp)
	client.SendToSelf(evt)

	rsp = &DeleteFriendResponse{FriendID: userID, RoomID: roomID}
	evt = newWebsocketEvent("delete-friend", rsp)
	s.hub.BroadcastToUsers(evt, friendID)

	return nil
}
