package ws

import (
	"context"
	"encoding/json"
	"gochat/src/core"
	"time"
)

func (s *Server) handleEvent(ctx context.Context, client *client, jsonMessage []byte) {
	var data json.RawMessage
	msg := WebsocketEvent{Data: &data}
	if err := json.Unmarshal(jsonMessage, &msg); err != nil {
		return
	}

	switch msg.Action {
	case "initialize":
		var req InitializeRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.initialize(ctx, client, &req)
	case "send-message":
		var req SendMessageRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.sendMessage(ctx, client, &req)
	case "add-friend":
		var req AddFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.addFriend(ctx, client, &req)
	case "accept-friend":
		var req AcceptFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.acceptFriend(ctx, client, &req)
	case "refuse-friend":
		var req RefuseFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.refuseFriend(ctx, client, &req)
	case "delete-friend":
		var req DeleteFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.deleteFriend(ctx, client, &req)
	case "new-room":
		var req NewRoomRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.newRoom(ctx, client, &req)
	case "update-room":
		var req UpdateRoomResponse
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.updateRoom(ctx, client, &req)
	case "delete-room":
		var req DeleteRoomRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.deleteRoom(ctx, client, &req)
	case "leave-room":
		var req LeaveRoomRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.leaveRoom(ctx, client, &req)
	case "add-members":
		var req AddMembersRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.addMembers(ctx, client, &req)
	case "delete-members":
		var req DeleteMembersRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		s.deleteMembers(ctx, client, &req)
	}
}

func (s *Server) initialize(ctx context.Context, client *client, req *InitializeRequest) {
	rooms, err := s.GetUserRooms(ctx, client.userID)
	if err != nil {
		client.sendMsg("toast", "internal server error")
		return
	}

	if err = s.GetRoomMessages(ctx, rooms); err != nil {
		return
	}

	friends, err := s.GetUserFriends(ctx, client.userID)
	if err != nil {
		client.sendMsg("toast", "internal server error")
		return
	}

	s.hub.registerClient(client, rooms)

	rsp := &InitializeResponse{
		Rooms:   rooms,
		Friends: friends,
	}
	client.sendMsg("initialize", rsp)
}

func (s *Server) sendMessage(ctx context.Context, client *client, req *SendMessageRequest) {
	myHubs := s.hub.userHubs[client.userID]
	if _, ok := myHubs[req.RoomID]; !ok {
		client.sendMsg("toast", "invalid room id")
		return
	}

	user, err := s.GetUser(ctx, client.userID)
	if err != nil {
		return
	}

	message := &core.MessageInfo{
		RoomID:   req.RoomID,
		SenderID: user.ID,
		Name:     user.Nickname,
		Avatar:   user.Avatar,
		Content:  req.Content,
		Kind:     req.Kind,
		Divide:   false,
		SendAt:   time.Now(),
	}
	if err = s.CacheMessage(ctx, req.RoomID, message); err != nil {
		return
	}

	rsp := &SendMessageResponse{Message: message}
	s.hub.broadcastByRoom("send-message", rsp, req.RoomID)
}

func (s *Server) addFriend(ctx context.Context, client *client, req *AddFriendRequest) {
}
func (s *Server) acceptFriend(ctx context.Context, client *client, req *AcceptFriendRequest) {
}
func (s *Server) refuseFriend(ctx context.Context, client *client, req *RefuseFriendRequest) {
}
func (s *Server) deleteFriend(ctx context.Context, client *client, req *DeleteFriendRequest) {
}
func (s *Server) newRoom(ctx context.Context, client *client, req *NewRoomRequest) {
}
func (s *Server) updateRoom(ctx context.Context, client *client, req *UpdateRoomResponse) {
}
func (s *Server) deleteRoom(ctx context.Context, client *client, req *DeleteRoomRequest) {
}
func (s *Server) leaveRoom(ctx context.Context, client *client, req *LeaveRoomRequest) {
}
func (s *Server) addMembers(ctx context.Context, client *client, req *AddMembersRequest) {
}
func (s *Server) deleteMembers(ctx context.Context, client *client, req *DeleteMembersRequest) {
}
