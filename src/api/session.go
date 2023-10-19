package api

import (
	db "gochat/src/db/sqlc"
	"gochat/src/util/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id"`
	ClientIp  string    `json:"client_ip"`
	UserAgent string    `json:"user_agent"`
	ExpireAt  time.Time `json:"expire_at"`
	CreateAt  time.Time `json:"create_at"`
}

// ======================== // deleteSession // ======================== //

type DeleteSessionRequest struct {
	SessionID uuid.UUID `uri:"session_id" binding:"required,uuid4"`
}

func (s *Server) deleteSession(ctx *gin.Context) {
	var req DeleteSessionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	err := s.store.RemoveSession(ctx, req.SessionID)
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}
}

// ======================== // ListSessions // ======================== //

type ListSessionsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}
type ListSessionsResponse struct {
	Total    int64      `json:"total"`
	Sessions []*Session `json:"sessions"`
}

func (s *Server) listSessions(ctx *gin.Context) {
	var req ListSessionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	sessions, err := s.store.ListSessions(ctx, &db.ListSessionsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		UserID: authPayload.UserID,
	})
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}

	ctx.JSON(http.StatusOK, convertListSessions(sessions))
}

func convertListSessions(sessions []*db.ListSessionsRow) *ListSessionsResponse {
	if len(sessions) == 0 {
		return &ListSessionsResponse{}
	}

	sess := make([]*Session, 0, 5)
	for _, session := range sessions {
		sess = append(sess, &Session{
			ID:        session.ID,
			ClientIp:  session.ClientIp,
			UserAgent: session.UserAgent,
			ExpireAt:  session.ExpireAt,
			CreateAt:  session.CreateAt,
		})
	}

	return &ListSessionsResponse{
		Total:    sessions[0].Total,
		Sessions: sess,
	}
}
