package hub

type Room struct {
	clients   map[*Client]bool
	join      chan *Client
	leave     chan *leaveRoomParams
	broadcast chan []byte
}

func startNewRoom() *Room {
	room := &Room{
		clients:   make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *leaveRoomParams),
		broadcast: make(chan []byte, 128),
	}
	go room.run()
	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case arg := <-r.leave:
			delete(r.clients, arg.client)
			if len(r.clients) == 0 {
				delete(arg.hub.rooms, arg.roomID)
				if userHubs, ok := arg.hub.userRooms[arg.client.userID]; ok {
					delete(userHubs, arg.roomID)
				}
				close(r.broadcast)
			}
			arg.done <- true
		case message, ok := <-r.broadcast:
			if !ok {
				return
			}
			for client := range r.clients {
				client.send <- message
			}
		}
	}
}

type leaveRoomParams struct {
	roomID int64
	client *Client
	hub    *HubCenter
	done   chan bool
}

func newLeaveRoomParams(roomID int64, client *Client, hub *HubCenter) *leaveRoomParams {
	return &leaveRoomParams{
		roomID: roomID,
		client: client,
		hub:    hub,
		done:   make(chan bool),
	}
}
