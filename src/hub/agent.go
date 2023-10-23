package hub

import (
	"encoding/json"
	"gochat/src/db"
)

func (h *Hub) Register(client *Client, rooms []*db.RoomInfo) {
	if h.running {
		arg := &registerParams{client: client, rooms: rooms}
		h.register <- arg
	}
}

func (h *Hub) Unregister(client *Client) {
	if h.running {
		h.unregister <- client
	}
}

func (h *Hub) JoinRoom(roomID int64, userIDs ...int64) {
	if h.running {
		arg := &roomParams{roomID: roomID, userIDs: userIDs}
		h.join <- arg
	}
}

func (h *Hub) LeaveRoom(roomID int64, userIDs ...int64) {
	if h.running {
		arg := &roomParams{roomID: roomID, userIDs: userIDs}
		h.join <- arg
	}
}

func (h *Hub) DeleteRoom(roomID int64) {
	if h.running {
		h.delete <- roomID
	}
}

func (h Hub) IsUserInRoom(userID int64, roomID int64) bool {
	if roomMap, ok := h.userRooms[userID]; ok {
		if _, ok = roomMap[roomID]; ok {
			return true
		}
	}
	return false
}

func (h *Hub) BroadcastToRoom(data any, roomID int64) {
	message, err := json.Marshal(data)
	if err != nil {
		return
	}
	if room, ok := h.rooms[roomID]; ok {
		room.broadcast <- message
	}
}

func (h *Hub) BroadcastToUsers(data any, userIDs ...int64) {
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
