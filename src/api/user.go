package api

import (
	"gochat/src/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

// ======================== // CreateUser // ======================== //

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,alphanum"`
	Role     string `json:"role" binding:"required,alphanum"`
}

type CreateUserResponse struct {
	User User `json:"user"`
}

func (s *Server) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	res, err := s.store.CreateUserTx(c, db.CreateUserTxParams{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		return
	}

	rsp := &CreateUserResponse{
		User: convertUser(&res.User),
	}

	c.JSON(http.StatusOK, rsp)
}

// ======================== // DeleteUser // ======================== //

type DeleteUserRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func (s *Server) DeleteUser(c *gin.Context) {

}

// ======================== // ListUsers // ======================== //

type ListUsersRequest struct {
	PageID   int64 `json:"page_id" binding:"required"`
	PageSize int64 `json:"page_size" binding:"required"`
}

type ListUsersResponse struct {
	Total int64  `json:"total"`
	Users []User `json:"users"`
}

func (s *Server) ListUsers(c *gin.Context) {

}

// ======================== // UpdateUser // ======================== //

type UpdateUserRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,alphanum"`
	Nickname string `json:"nickname" binding:"required,alphanum"`
	Avatar   string `json:"avatar" binding:"required,alphanum"`
	Role     string `json:"role" binding:"required,alphanum"`
	Deleted  string `json:"deleted" binding:"required,alphanum"`
}

func (s *Server) UpdateUser(ctx *gin.Context) {

}
