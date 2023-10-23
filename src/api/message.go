package api

import (
	"context"
	"gochat/src/db"
	"gochat/src/hub"
	"time"
)

// ======================== // initialize // ======================== //

type InitializeRequest struct {
	EndTime time.Time `json:"end_time"`
}
type InitializeResponse struct {
	Rooms   []*db.RoomInfo   `json:"rooms"`
	Friends []*db.FriendInfo `json:"friends"`
}

func (s *Server) initialize(ctx context.Context, client *hub.Client, req *InitializeRequest) error {
	userID := client.GetUserID()

	rooms, err := s.store.GetUserRooms(ctx, userID, req.EndTime)
	if err != nil {
		return errInternal
	}

	friends, err := s.store.GetUserFriends(ctx, userID)
	if err != nil {
		return errInternal
	}

	s.hub.Register(client, rooms)

	rsp := &InitializeResponse{
		Rooms:   rooms,
		Friends: friends,
	}
	evt := newWebsocketEvent("initialize", rsp)
	client.SendToSelf(evt)

	return nil
}

// ======================== // sendMessage // ======================== //

type SendMessageRequest struct {
	RoomID  int64  `json:"room_id"`
	Content string `json:"content"`
	Kind    string `json:"kind"`
}
type SendMessageResponse struct {
	Message *db.MessageInfo `json:"message"`
}

func (s *Server) sendMessage(ctx context.Context, client *hub.Client, req *SendMessageRequest) error {
	userID := client.GetUserID()

	if ok := s.hub.IsUserInRoom(userID, req.RoomID); !ok {
		return nil
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return errInternal
	}

	message := &db.MessageInfo{
		RoomID:   req.RoomID,
		SenderID: user.ID,
		Name:     user.Nickname,
		Avatar:   user.Avatar,
		Content:  req.Content,
		Kind:     req.Kind,
		Divide:   false,
		SendAt:   time.Now(),
	}
	if err = s.store.CacheMessage(ctx, message); err != nil {
		return errInternal
	}

	rsp := &SendMessageResponse{Message: message}
	evt := newWebsocketEvent("send-message", rsp)
	s.hub.BroadcastToRoom(evt, req.RoomID)

	return nil
}
