package ws

import (
	"gochat/src/core"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*core.State
	hub *hubService
}

func NewServer(state *core.State) *Server {
	hub := newHubService()
	go hub.run()

	return &Server{
		State: state,
		hub:   hub,
	}
}

func (s *Server) RegisterRouter(router *gin.Engine) {
	router.GET("/ws/:authorization", s.handleWebsocket)
}

// ======================== // handleWebsocket // ======================== //

type WebSocketRequest struct {
	Authorization string `uri:"authorization" binding:"required"`
}

func (s *Server) handleWebsocket(ctx *gin.Context) {
	var req WebSocketRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		return
	}

	payload, err := s.TokenMaker.VerifyToken(req.Authorization)
	if err != nil {
		return
	}

	user, err := s.GetUser(ctx, payload.UserID)
	if err != nil {
		return
	}

	s.serveWs(ctx, user)
}
