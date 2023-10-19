package api

import (
	"context"
	"encoding/json"
	"gochat/src/hub"
)

type WebsocketEvent struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

func newWebsocketEvent(action string, data interface{}) *WebsocketEvent {
	return &WebsocketEvent{Action: action, Data: data}
}

func (s *Server) HandleEvent(ctx context.Context, client *hub.Client, jsonMessage []byte) {
	var data json.RawMessage
	msg := WebsocketEvent{Data: &data}
	err := json.Unmarshal(jsonMessage, &msg)
	if err != nil {
		return
	}

	switch msg.Action {
	case "initialize":
		var req InitializeRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.initialize(ctx, client, &req)
	case "send-message":
		var req SendMessageRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.sendMessage(ctx, client, &req)
	case "add-friend":
		var req AddFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.addFriend(ctx, client, &req)
	case "accept-friend":
		var req AcceptFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.acceptFriend(ctx, client, &req)
	case "refuse-friend":
		var req RefuseFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.refuseFriend(ctx, client, &req)
	case "delete-friend":
		var req DeleteFriendRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.deleteFriend(ctx, client, &req)
	case "new-room":
		var req NewRoomRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.newRoom(ctx, client, &req)
	case "update-room":
		var req UpdateRoomResponse
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.updateRoom(ctx, client, &req)
	case "delete-room":
		var req DeleteRoomRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.deleteRoom(ctx, client, &req)
	case "leave-room":
		var req LeaveRoomRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.leaveRoom(ctx, client, &req)
	case "add-members":
		var req AddMembersRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.addMembers(ctx, client, &req)
	case "delete-members":
		var req DeleteMembersRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		err = s.deleteMembers(ctx, client, &req)
	}

	if err != nil {
		evt := newWebsocketEvent("toast", err.Error())
		client.SendToSelf(evt)
	}
}

func (s *Server) Unregister(client *hub.Client) {
	s.hub.Unregister(client)
}
