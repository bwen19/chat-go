package ws

import (
	db "gochat/src/db/sqlc"
	"log"
)

// ======================== // hub // ======================== //

type hub struct {
	clients   map[*client]bool
	join      chan *client
	left      chan *client
	broadcast chan []byte
	close     chan bool
}

func runNewHub() *hub {
	hub := &hub{
		clients:   make(map[*client]bool),
		join:      make(chan *client),
		left:      make(chan *client),
		broadcast: make(chan []byte, 128),
		close:     make(chan bool),
	}
	go hub.run()
	return hub
}

func (h *hub) run() {
	defer func() {
		log.Println("exit hub")
	}()
	for {
		select {
		case client := <-h.join:
			h.clients[client] = true
		case client := <-h.left:
			delete(h.clients, client)
		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		case <-h.close:
			return
		}
	}
}

// ======================== // hubService // ======================== //

type hubService struct {
	hubs       map[int64]*hub           // roomID - hub
	privateHub map[int64]int64          // userID - roomID
	userHubs   map[int64]map[int64]bool // userID - roomID - bool
	register   chan *hubRegisterParams
	unregister chan *hubUnregisterParams
	addHub     chan int64
	delHub     chan int64
}

func newHubService() *hubService {
	return &hubService{
		hubs:       make(map[int64]*hub),
		privateHub: make(map[int64]int64),
		userHubs:   make(map[int64]map[int64]bool),
		register:   make(chan *hubRegisterParams),
		unregister: make(chan *hubUnregisterParams),
		addHub:     make(chan int64),
		delHub:     make(chan int64),
	}
}

func (h *hubService) run() {
	for {
		select {
		case params := <-h.register:
			if _, ok := h.privateHub[params.userID]; !ok {
				h.privateHub[params.userID] = params.privateRoom
			}
			if _, ok := h.userHubs[params.userID]; !ok {
				roomSet := make(map[int64]bool, len(params.roomIDs))
				for _, roomID := range params.roomIDs {
					roomSet[roomID] = true
					if _, ok := h.hubs[roomID]; !ok {
						h.hubs[roomID] = runNewHub()
					}
				}
				h.userHubs[params.userID] = roomSet
			}
			for _, roomID := range params.roomIDs {
				h.hubs[roomID].join <- params.client
			}
		case params := <-h.unregister:
			for roomID := range h.userHubs[params.userID] {
				if hub, ok := h.hubs[roomID]; ok {
					hub.left <- params.client
					if len(hub.clients) == 0 {
						hub.close <- true
						delete(h.hubs, roomID)
					}
				}
			}
			if privateRoomID, ok := h.privateHub[params.userID]; ok {
				if _, ok := h.hubs[privateRoomID]; !ok {
					delete(h.privateHub, params.userID)
					delete(h.userHubs, params.userID)
				}
			}
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
	userID      int64
	privateRoom int64
	roomIDs     []int64
	client      *client
}

func (h *hubService) registerClient(user *db.User, client *client) {
	arg := &hubRegisterParams{
		userID:      user.ID,
		privateRoom: user.RoomID,
		client:      client,
	}
	h.register <- arg
}

type hubUnregisterParams struct {
	userID int64
	client *client
}

func (h *hubService) unregisterClient(userID int64, client *client) {
	arg := &hubUnregisterParams{
		userID: userID,
		client: client,
	}
	h.unregister <- arg
}
