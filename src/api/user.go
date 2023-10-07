package api

import (
	"errors"
	db "gochat/src/db/sqlc"
	"gochat/src/util"
	"gochat/src/util/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Nickname string    `json:"nickname"`
	Role     string    `json:"role"`
	RoomID   int64     `json:"room_id"`
	Deleted  bool      `json:"deleted"`
	CreateAt time.Time `json:"create_at"`
}

// ======================== // createUser // ======================== //

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}
type CreateUserResponse struct {
	User *User `json:"user"`
}

func (s *Server) createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := s.Store.CreateUserTx(c, db.CreateUserTxParams{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			c.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := &CreateUserResponse{
		User: convertUser(&res.User),
	}
	c.JSON(http.StatusOK, rsp)
}

// ======================== // deleteUser // ======================== //

type DeleteUserRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

func (s *Server) deleteUser(c *gin.Context) {
	var req DeleteUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.Store.DeleteUserTx(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

// ======================== // listUsers // ======================== //

type ListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type ListUsersResponse struct {
	Total int64   `json:"total"`
	Users []*User `json:"users"`
}

func (s *Server) listUsers(c *gin.Context) {
	var req ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := s.Store.ListUsers(c, db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, convertListUsers(users))
}

// ======================== // updateUser // ======================== //

type UpdateUserRequest struct {
	UserID   int64   `json:"user_id" binding:"required"`
	Username *string `json:"username" binding:"omitempty,alphanum,min=2,max=50"`
	Password *string `json:"password" binding:"omitempty,min=6,max=50"`
	Nickname *string `json:"nickname" binding:"omitempty,min=2,max=50"`
	Avatar   *string `json:"avatar" binding:"-"`
	Role     *string `json:"role" binding:"omitempty,oneof=admin user"`
	Deleted  *bool   `json:"deleted" binding:"-"`
}
type UpdateUserResponse struct {
	User *User `json:"user"`
}

func (s *Server) updateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{ID: req.UserID}
	if req.Username != nil {
		arg.Username = pgtype.Text{String: *req.Username, Valid: true}
	}
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		arg.HashedPassword = pgtype.Text{String: hashedPassword, Valid: true}
	}
	if req.Nickname != nil {
		arg.Nickname = pgtype.Text{String: *req.Nickname, Valid: true}
	}
	if req.Avatar != nil {
		arg.Avatar = pgtype.Text{String: *req.Avatar, Valid: true}
	}
	if req.Role != nil {
		arg.Role = pgtype.Text{String: *req.Role, Valid: true}
	}
	if req.Deleted != nil {
		arg.Deleted = pgtype.Bool{Bool: *req.Deleted, Valid: true}
	}

	user, err := s.Store.UpdateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := &UpdateUserResponse{User: convertUser(&user)}
	c.JSON(http.StatusOK, rsp)
}

// ======================== // changePassword // ======================== //

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=2,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=2,max=50"`
}

func (s *Server) changePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := s.Store.GetUser(c, payload.UserID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = util.CheckPassword(req.OldPassword, user.HashedPassword); err != nil {
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	_, err = s.Store.UpdateUser(c, db.UpdateUserParams{
		ID:             payload.UserID,
		HashedPassword: pgtype.Text{String: hashedPassword, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
}

// ======================== // changeUserInfo // ======================== //

type ChangeUserInfoRequest struct {
	Username *string `json:"username" binding:"omitempty,alphanum,min=2,max=50"`
	Nickname *string `json:"nickname" binding:"omitempty,min=2,max=50"`
	Avatar   *string `json:"avatar" binding:"-"`
}
type ChangeUserInfoResponse struct {
	User *User `json:"user"`
}

func (s *Server) changeUserInfo(c *gin.Context) {
	var req ChangeUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.UpdateUserParams{ID: authPayload.UserID}
	if req.Username != nil {
		arg.Username = pgtype.Text{String: *req.Username, Valid: true}
	}
	if req.Nickname != nil {
		arg.Nickname = pgtype.Text{String: *req.Nickname, Valid: true}
	}
	if req.Avatar != nil {
		arg.Avatar = pgtype.Text{String: *req.Avatar, Valid: true}
	}

	user, err := s.Store.UpdateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := &UpdateUserResponse{User: convertUser(&user)}
	c.JSON(http.StatusOK, rsp)
}
