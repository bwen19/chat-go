package hub

import (
	"gochat/src/db"
)

type HubCenter struct {
	rooms      map[int64]*Room          // roomID - room
	pRoom      map[int64]int64          // userID - personal roomID
	userRooms  map[int64]map[int64]bool // userID - roomID - bool
	running    bool
	register   chan *registerParams
	unregister chan *Client
	join       chan *roomParams
	leave      chan *roomParams
}

func NewHubCenter() *HubCenter {
	return &HubCenter{
		rooms:      make(map[int64]*Room),
		pRoom:      make(map[int64]int64),
		userRooms:  make(map[int64]map[int64]bool),
		register:   make(chan *registerParams),
		unregister: make(chan *Client),
		join:       make(chan *roomParams),
		leave:      make(chan *roomParams),
	}
}

func (h *HubCenter) Start() {
	if h.running {
		return
	}
	h.running = true
	go h.run()
}

func (h *HubCenter) run() {
	for {
		select {
		case arg := <-h.register:
			h.registerClient(arg)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case arg := <-h.join:
			h.joinRoom(arg)
		case arg := <-h.leave:
			h.leaveRoom(arg)
		}
	}
}

type registerParams struct {
	client *Client
	rooms  []*db.RoomInfo
}

func (h *HubCenter) registerClient(arg *registerParams) {
	userID := arg.client.userID
	privateRoomID := arg.client.roomID

	if _, ok := h.pRoom[arg.client.userID]; !ok {
		h.pRoom[arg.client.userID] = privateRoomID
	}

	if _, ok := h.userRooms[userID]; !ok {
		roomSet := make(map[int64]bool, len(arg.rooms))
		for _, room := range arg.rooms {
			roomID := room.ID
			if _, ok := h.rooms[roomID]; !ok {
				h.rooms[roomID] = startNewRoom()
			}
			roomSet[roomID] = true
		}
		h.userRooms[userID] = roomSet
	}

	for roomID := range h.userRooms[userID] {
		h.rooms[roomID].join <- arg.client
	}
}

func (h *HubCenter) unregisterClient(client *Client) {
	userID := client.userID

	for roomID := range h.userRooms[userID] {
		if hub, ok := h.rooms[roomID]; ok {
			arg := newLeaveRoomParams(roomID, client, h)
			hub.leave <- arg
			<-arg.done
		}
	}
	if privateRoomID, ok := h.pRoom[userID]; ok {
		if _, ok := h.rooms[privateRoomID]; !ok {
			delete(h.pRoom, userID)
			delete(h.userRooms, userID)
		}
	}
	close(client.send)
}

type roomParams struct {
	roomID  int64
	userIDs []int64
}

func (h *HubCenter) joinRoom(arg *roomParams) {
	roomID := arg.roomID
	hub, ok := h.rooms[roomID]
	if !ok {
		hub = startNewRoom()
		h.rooms[roomID] = hub
	}

	for _, userID := range arg.userIDs {
		if privateRoomID, ok := h.pRoom[userID]; ok {
			if pHub, ok := h.rooms[privateRoomID]; ok {
				for client := range pHub.clients {
					hub.join <- client
				}
			}
		}
	}
}

func (h *HubCenter) leaveRoom(arg *roomParams) {
	roomID := arg.roomID
	if hub, ok := h.rooms[roomID]; ok {
		for _, userID := range arg.userIDs {
			if privateRoomID, ok := h.pRoom[userID]; ok {
				if pHub, ok := h.rooms[privateRoomID]; ok {
					for client := range pHub.clients {
						lvarg := newLeaveRoomParams(roomID, client, h)
						hub.leave <- lvarg
						<-lvarg.done
					}
				}
			}
		}
	}
}
