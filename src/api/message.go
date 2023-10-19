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

	rooms, err := s.store.GetUserRooms(ctx, userID)
	if err != nil {
		return ErrInternal
	}

	if err = s.store.GetRoomMessages(ctx, rooms, req.EndTime); err != nil {
		return ErrInternal
	}

	friends, err := s.store.GetUserFriends(ctx, userID)
	if err != nil {
		return ErrInternal
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
	if ok := s.hub.IsUserInRoom(client, req.RoomID); !ok {
		return nil
	}

	user, err := s.store.GetUserByID(ctx, client.GetUserID())
	if err != nil {
		return ErrInternal
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
		return ErrInternal
	}

	rsp := &SendMessageResponse{Message: message}
	env := newWebsocketEvent("initialize", rsp)
	s.hub.BroadcastToRoom(env, req.RoomID)

	return nil
}
