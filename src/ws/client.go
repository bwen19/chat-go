package ws

import (
	"encoding/json"
	"gochat/src/core"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type client struct {
	userID int64
	roomID int64
	conn   *websocket.Conn
	send   chan []byte
}

func (c *client) readPump(ctx *gin.Context, srv *Server) {
	defer func() {
		srv.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if msgType == websocket.TextMessage {
			srv.handleEvent(ctx, c, message)
		} else if msgType == websocket.CloseMessage {
			break
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *client) sendMsg(action string, data any) {
	wsEvent := WebsocketEvent{Action: action, Data: data}
	message, err := json.Marshal(wsEvent)
	if err != nil {
		return
	}
	c.send <- message
}

// ======================== // serveWs // ======================== //

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) serveWs(ctx *gin.Context, user *core.UserInfo) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := &client{
		userID: user.ID,
		roomID: user.RoomID,
		conn:   conn,
		send:   make(chan []byte, 128),
	}

	go client.readPump(ctx, s)
	go client.writePump()
}
