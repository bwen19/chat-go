package hub

import (
	"encoding/json"
	"gochat/src/db"
)

func (h *HubCenter) Register(client *Client, rooms []*db.RoomInfo) {
	arg := &registerParams{
		client: client,
		rooms:  rooms,
	}
	h.register <- arg
}

func (h *HubCenter) Unregister(client *Client) {
	h.unregister <- client
}

func (h *HubCenter) JoinRoom(roomID int64, userIDs ...int64) {
	arg := &roomParams{roomID: roomID, userIDs: userIDs}
	h.join <- arg
}

func (h *HubCenter) LeaveRoom(roomID int64, userIDs ...int64) {
	arg := &roomParams{roomID: roomID, userIDs: userIDs}
	h.join <- arg
}

func (h HubCenter) IsUserInRoom(client *Client, roomID int64) bool {
	myHubs := h.userRooms[client.userID]
	_, ok := myHubs[roomID]
	return ok
}

func (h *HubCenter) BroadcastToRoom(data any, roomID int64) {
	message, err := json.Marshal(data)
	if err != nil {
		return
	}
	h.rooms[roomID].broadcast <- message
}

func (h *HubCenter) BroadcastToUsers(data any, userIDs ...int64) {
	message, err := json.Marshal(data)
	if err != nil {
		return
	}

	for _, userID := range userIDs {
		if roomID, ok := h.pRoom[userID]; ok {
			h.rooms[roomID].broadcast <- message
		}
	}
}
