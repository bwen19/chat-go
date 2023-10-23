package hub

import (
	"gochat/src/db"
)

type Hub struct {
	running    bool
	rooms      map[int64]*Room              // roomID - room
	pRoom      map[int64]int64              // userID - personal roomID
	userRooms  map[int64]map[int64]struct{} // userID - roomIDs
	register   chan *registerParams
	unregister chan *Client
	join       chan *roomParams
	leave      chan *roomParams
	delete     chan int64
	close      chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[int64]*Room),
		pRoom:      make(map[int64]int64),
		userRooms:  make(map[int64]map[int64]struct{}),
		register:   make(chan *registerParams),
		unregister: make(chan *Client),
		join:       make(chan *roomParams),
		leave:      make(chan *roomParams),
		delete:     make(chan int64),
		close:      make(chan struct{}),
	}
}

func (h *Hub) Start() {
	if h.running {
		return
	}
	h.running = true
	go h.run()
}

func (h *Hub) Stop() {
	if h.running {
		h.close <- struct{}{}
		h.running = false
	}
}

type registerParams struct {
	client *Client
	rooms  []*db.RoomInfo
}

type roomParams struct {
	roomID  int64
	userIDs []int64
}

func (h *Hub) run() {
	for {
		select {
		case arg := <-h.register:
			userID := arg.client.userID
			pRoomID := arg.client.roomID

			if _, ok := h.pRoom[userID]; !ok {
				h.pRoom[userID] = pRoomID
			}

			if _, ok := h.userRooms[userID]; !ok {
				roomMap := make(map[int64]struct{}, len(arg.rooms))
				for _, room := range arg.rooms {
					roomID := room.ID
					if _, ok := h.rooms[roomID]; !ok {
						h.rooms[roomID] = startNewRoom()
					}
					roomMap[roomID] = struct{}{}
				}
				h.userRooms[userID] = roomMap
			}

			for roomID := range h.userRooms[userID] {
				h.rooms[roomID].join <- arg.client
			}

		case client := <-h.unregister:
			if roomMap, ok := h.userRooms[client.userID]; ok {
				for roomID := range roomMap {
					if room, ok := h.rooms[roomID]; ok {
						room.leave <- client
					}
				}
			}
			close(client.send)

		case arg := <-h.join:
			roomID := arg.roomID

			room, ok := h.rooms[roomID]
			if !ok {
				room = startNewRoom()
				h.rooms[roomID] = room
			}

			for _, userID := range arg.userIDs {
				if pRoomID, ok := h.pRoom[userID]; ok {
					if pRoom, ok := h.rooms[pRoomID]; ok {
						for client := range pRoom.clients {
							room.join <- client
						}
					}
				}
				if roomMap, ok := h.userRooms[userID]; ok {
					roomMap[roomID] = struct{}{}
				}
			}

		case arg := <-h.leave:
			roomID := arg.roomID
			if room, ok := h.rooms[roomID]; ok {
				for _, userID := range arg.userIDs {
					if roomMap, ok := h.userRooms[userID]; ok {
						delete(roomMap, roomID)
					}
					if pRoomID, ok := h.pRoom[userID]; ok {
						if pRoom, ok := h.rooms[pRoomID]; ok {
							for client := range pRoom.clients {
								room.leave <- client
							}
						}
					}
				}
			}

		case roomID := <-h.delete:
			if room, ok := h.rooms[roomID]; ok {
				delete(h.rooms, roomID)
				var userID int64
				for client := range room.clients {
					if client.userID != userID {
						userID = client.userID
						if roomMap, ok := h.userRooms[userID]; ok {
							delete(roomMap, roomID)
						}
					}
				}
				close(room.broadcast)
			}

		case <-h.close:
			return
		}
	}
}
