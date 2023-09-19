package api

import (
	"context"
	"errors"
	"gochat/src/db"
	"gochat/src/pb"
	"gochat/src/utils"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if violations := validateLoginRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	if err = utils.CheckPassword(req.GetPassword(), user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	mtdt := server.extractLoginInfo(ctx)
	err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ClientIp:     mtdt.ClientIp,
		UserAgent:    mtdt.UserAgent,
		ExpireAt:     refreshPayload.ExpireAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rsp := &pb.LoginResponse{
		User:         convertUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return rsp, nil
}

func validateLoginRequest(req *pb.LoginRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateName(req.GetUsername(), 3, 50); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := utils.ValidateString(req.GetPassword(), 6, 50); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
