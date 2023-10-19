package api

import (
	"gochat/src/db"
	"gochat/src/db/sqlc"
	"gochat/src/util"
	"gochat/src/util/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// ======================== // createUser // ======================== //

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=2,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}
type CreateUserResponse struct {
	User *db.UserInfo `json:"user"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	user, err := s.store.CreateUser(ctx, &db.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		UniqueViolationResponse(ctx, err)
		return
	}

	rsp := &CreateUserResponse{User: user}
	ctx.JSON(http.StatusOK, rsp)
}

// ======================== // deleteUser // ======================== //

type DeleteUserRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

func (s *Server) deleteUser(ctx *gin.Context) {
	var req DeleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	err := s.store.RemoveUser(ctx, req.UserID)
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

// ======================== // listUsers // ======================== //

type ListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type ListUsersResponse struct {
	Total int64          `json:"total"`
	Users []*db.UserInfo `json:"users"`
}

func (s *Server) listUsers(ctx *gin.Context) {
	var req ListUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	total, users, err := s.store.GetUsers(ctx, &sqlc.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		InternalErrorResponse(ctx)
		return
	}

	rsp := &ListUsersResponse{Total: total, Users: users}
	ctx.JSON(http.StatusOK, rsp)
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
	User *db.UserInfo `json:"user"`
}

func (s *Server) updateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	arg := &sqlc.UpdateUserParams{ID: req.UserID}
	if req.Username != nil {
		arg.Username = pgtype.Text{String: *req.Username, Valid: true}
	}
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(*req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
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

	user, err := s.store.ModifyUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := &UpdateUserResponse{User: user}
	ctx.JSON(http.StatusOK, rsp)
}

// ======================== // changePassword // ======================== //

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=2,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=2,max=50"`
}

func (s *Server) changePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := s.store.RetrieveUserByID(ctx, payload.UserID)
	if err != nil {
		RecordNotFoundResponse(ctx, err)
		return
	}

	if err = util.CheckPassword(req.OldPassword, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	_, err = s.store.UpdateUser(ctx, &sqlc.UpdateUserParams{
		ID:             payload.UserID,
		HashedPassword: pgtype.Text{String: hashedPassword, Valid: true},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
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
	User *db.UserInfo `json:"user"`
}

func (s *Server) changeUserInfo(ctx *gin.Context) {
	var req ChangeUserInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := &sqlc.UpdateUserParams{ID: authPayload.UserID}
	if req.Username != nil {
		arg.Username = pgtype.Text{String: *req.Username, Valid: true}
	}
	if req.Nickname != nil {
		arg.Nickname = pgtype.Text{String: *req.Nickname, Valid: true}
	}
	if req.Avatar != nil {
		arg.Avatar = pgtype.Text{String: *req.Avatar, Valid: true}
	}

	user, err := s.store.ModifyUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := &UpdateUserResponse{User: user}
	ctx.JSON(http.StatusOK, rsp)
}

// ======================== // FindUser // ======================== //

type FindUserRequest struct {
	Username string `form:"username" binding:"required,alphanum,min=2,max=50"`
}
type FindUserResponse struct {
	User *db.UserInfo `json:"user"`
}

func (s *Server) findUser(ctx *gin.Context) {
	var req FindUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		InvalidArgumentResponse(ctx)
		return
	}

	user, err := s.store.GetUserByName(ctx, req.Username)
	if err != nil {
		RecordNotFoundResponse(ctx, err)
		return
	}

	rsp := &FindUserResponse{User: db.NewUserInfo(user)}
	ctx.JSON(http.StatusOK, rsp)
}
