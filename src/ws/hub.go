package ws

import (
	"encoding/json"
	"gochat/src/core"
	"log"
)

// ======================== // hub // ======================== //

type hub struct {
	clients   map[*client]bool
	join      chan *client
	left      chan *client
	broadcast chan []byte
}

func runNewHub() *hub {
	hub := &hub{
		clients:   make(map[*client]bool),
		join:      make(chan *client),
		left:      make(chan *client),
		broadcast: make(chan []byte, 128),
	}
	go hub.run()
	return hub
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
		case client := <-h.left:
			delete(h.clients, client)
		case message, ok := <-h.broadcast:
			if !ok {
				return
			}
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}

// ======================== // hubService // ======================== //

type hubService struct {
	hubs       map[int64]*hub           // roomID - hub
	privateHub map[int64]int64          // userID - roomID
	userHubs   map[int64]map[int64]bool // userID - roomID - bool
	register   chan *hubRegisterParams
	unregister chan *client
	addHub     chan int64
	delHub     chan int64
}

func newHubService() *hubService {
	return &hubService{
		hubs:       make(map[int64]*hub),
		privateHub: make(map[int64]int64),
		userHubs:   make(map[int64]map[int64]bool),
		register:   make(chan *hubRegisterParams),
		unregister: make(chan *client),
		addHub:     make(chan int64),
		delHub:     make(chan int64),
	}
}

func (h *hubService) run() {
	log.Print("start Hub Service")
	for {
		select {
		case params := <-h.register:
			userID := params.client.userID
			privateRoomID := params.client.roomID

			if _, ok := h.privateHub[params.client.userID]; !ok {
				h.privateHub[params.client.userID] = privateRoomID
			}

			if _, ok := h.userHubs[userID]; !ok {
				roomSet := make(map[int64]bool, len(params.rooms))
				for _, room := range params.rooms {
					roomID := room.ID
					if _, ok := h.hubs[roomID]; !ok {
						h.hubs[roomID] = runNewHub()
					}
					roomSet[roomID] = true
				}
				h.userHubs[userID] = roomSet
			}

			for roomID := range h.userHubs[userID] {
				h.hubs[roomID].join <- params.client
			}
		case client := <-h.unregister:
			userID := client.userID

			for roomID := range h.userHubs[userID] {
				if hub, ok := h.hubs[roomID]; ok {
					if len(hub.clients) <= 1 {
						close(hub.broadcast)
						delete(h.hubs, roomID)
					} else {
						hub.left <- client
					}
				}
			}
			if privateRoomID, ok := h.privateHub[userID]; ok {
				if _, ok := h.hubs[privateRoomID]; !ok {
					delete(h.privateHub, userID)
					delete(h.userHubs, userID)
				}
			}
			close(client.send)
		case roomID := <-h.addHub:
			if _, ok := h.hubs[roomID]; !ok {
				h.hubs[roomID] = runNewHub()
			}
		case roomID := <-h.delHub:
			delete(h.hubs, roomID)
		}
	}
}

type hubRegisterParams struct {
	client *client
	rooms  []*core.RoomInfo
}

func (h *hubService) registerClient(client *client, rooms []*core.RoomInfo) {
	arg := &hubRegisterParams{
		client: client,
		rooms:  rooms,
	}
	h.register <- arg
}

func (h *hubService) broadcastByRoom(action string, data any, roomID int64) {
	wsEvent := WebsocketEvent{Action: action, Data: data}
	message, err := json.Marshal(wsEvent)
	if err != nil {
		return
	}

	h.hubs[roomID].broadcast <- message
}

func (h *hubService) broadcastByUsers(action string, data any, userIDs ...int64) {
	wsEvent := WebsocketEvent{Action: action, Data: data}
	message, err := json.Marshal(wsEvent)
	if err != nil {
		return
	}

	for _, userID := range userIDs {
		roomID := h.privateHub[userID]
		h.hubs[roomID].broadcast <- message
	}
}
