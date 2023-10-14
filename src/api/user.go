package api

import (
	"gochat/src/core"
	db "gochat/src/db/sqlc"
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
	User *core.UserInfo `json:"user"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.InvalidArgumentResponse(ctx)
		return
	}

	res, err := s.Store.CreateUserTx(ctx, db.CreateUserTxParams{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		util.UniqueViolationResponse(ctx, err)
		return
	}

	rsp := &CreateUserResponse{User: core.NewUserInfo(&res.User)}
	ctx.JSON(http.StatusOK, rsp)
}

// ======================== // deleteUser // ======================== //

type DeleteUserRequest struct {
	UserID int64 `uri:"user_id" binding:"required,min=1"`
}

func (s *Server) deleteUser(ctx *gin.Context) {
	var req DeleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		util.InvalidArgumentResponse(ctx)
		return
	}

	err := s.Store.DeleteUserTx(ctx, req.UserID)
	if err != nil {
		util.InternalErrorResponse(ctx)
		return
	}
}

// ======================== // listUsers // ======================== //

type ListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type ListUsersResponse struct {
	Total int64            `json:"total"`
	Users []*core.UserInfo `json:"users"`
}

func (s *Server) listUsers(ctx *gin.Context) {
	var req ListUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		util.InvalidArgumentResponse(ctx)
		return
	}

	users, err := s.Store.ListUsers(ctx, db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		util.InternalErrorResponse(ctx)
		return
	}

	ctx.JSON(http.StatusOK, convertListUsers(users))
}

func convertListUsers(users []db.ListUsersRow) *ListUsersResponse {
	if len(users) == 0 {
		return &ListUsersResponse{}
	}

	res := make([]*core.UserInfo, 0, 5)
	for _, user := range users {
		res = append(res, &core.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Deleted:  user.Deleted,
			CreateAt: user.CreateAt,
		})
	}

	return &ListUsersResponse{
		Total: users[0].Total,
		Users: res,
	}
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
	User *core.UserInfo `json:"user"`
}

func (s *Server) updateUser(ctx *gin.Context) {
	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.InvalidArgumentResponse(ctx)
		return
	}

	arg := db.UpdateUserParams{ID: req.UserID}
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

	user, err := s.Store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := &UpdateUserResponse{User: core.NewUserInfo(&user)}
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
		util.InvalidArgumentResponse(ctx)
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := s.Store.GetUser(ctx, payload.UserID)
	if err != nil {
		util.RecordNotFoundResponse(ctx, err)
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

	_, err = s.Store.UpdateUser(ctx, db.UpdateUserParams{
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
	User *core.UserInfo `json:"user"`
}

func (s *Server) changeUserInfo(ctx *gin.Context) {
	var req ChangeUserInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.InvalidArgumentResponse(ctx)
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

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

	user, err := s.Store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	rsp := &UpdateUserResponse{User: core.NewUserInfo(&user)}
	ctx.JSON(http.StatusOK, rsp)
}
