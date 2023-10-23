package hub

type Room struct {
	clients   map[*Client]struct{}
	join      chan *Client
	leave     chan *Client
	broadcast chan []byte
}

func startNewRoom() *Room {
	room := &Room{
		clients:   make(map[*Client]struct{}),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan []byte, 256),
	}
	go room.run()
	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}
		case client := <-r.leave:
			delete(r.clients, client)
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
