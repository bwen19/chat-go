package hub

import (
	"context"
	"encoding/json"
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

type ClientHandler interface {
	Unregister(client *Client)
	HandleEvent(ctx context.Context, client *Client, message []byte)
}

type Client struct {
	userID int64
	roomID int64
	conn   *websocket.Conn
	send   chan []byte
}

func NewClient(userID int64, roomID int64, conn *websocket.Conn) *Client {
	return &Client{
		userID: userID,
		roomID: roomID,
		conn:   conn,
		send:   make(chan []byte, 128),
	}
}

func (c *Client) ReadPump(ctx *gin.Context, h ClientHandler) {
	defer func() {
		h.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from Client
	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.TextMessage {
			h.HandleEvent(ctx, c, message)
		} else if msgType == websocket.CloseMessage {
			break
		}
	}
}

func (c *Client) WritePump() {
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

func (c *Client) GetUserID() int64 {
	return c.userID
}

func (c *Client) SendToSelf(data any) {
	message, err := json.Marshal(data)
	if err != nil {
		return
	}
	c.send <- message
}
