package ws

import (
	"errors"
	"gochat/src/util/state"
	"gochat/src/util/token"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	*state.State
	hub *hubService
}

func NewServer(state *state.State) *Server {
	hub := newHubService()
	log.Println("starting hub service")
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := s.verifyAuthorization(req.Authorization)
	if err != nil {
		if errors.Is(err, token.ErrExpiredToken) {
			ctx.JSON(http.StatusPaymentRequired, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := s.GetUser(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusNotFound, "")
		return
	}

	client := newClient(conn, payload)
	s.hub.registerClient(&user, client)

	go client.readPump(s, user.ID)
	go client.writePump()
}

func (s *Server) verifyAuthorization(authorization string) (*token.Payload, error) {
	fields := strings.Fields(authorization)
	if len(fields) < 2 {
		err := errors.New("invalid authorization format")
		return nil, err
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != "bearer" {
		err := errors.New("unsupported authorization type")
		return nil, err
	}

	accessToken := fields[1]
	payload, err := s.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
