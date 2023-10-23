package api

import (
	"gochat/src/hub"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketRequest struct {
	Authorization string `uri:"authorization" binding:"required"`
}

func (s *Server) handleWebsocket(ctx *gin.Context) {
	var req WebSocketRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		return
	}

	payload, err := s.tokenMaker.VerifyToken(req.Authorization)
	if err != nil {
		return
	}

	session, err := s.store.GetSession(ctx, payload.ID)
	if err != nil {
		return
	}

	if session.UserID != payload.UserID || session.RefreshToken != req.Authorization {
		return
	}

	user, err := s.store.GetUserByID(ctx, payload.UserID)
	if err != nil {
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	client := hub.NewClient(user.ID, user.RoomID, conn)

	go client.ReadPump(ctx, s)
	go client.WritePump()
}
