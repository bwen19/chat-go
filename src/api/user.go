package api

import (
	"context"
	"gochat/src/db"
	"gochat/src/pb"
	"gochat/src/utils"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	res, err := server.store.CreateUserTx(ctx, db.CreateUserTxParams{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(res.User),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := utils.ValidateName(req.GetUsername(), 3, 50); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := utils.ValidateString(req.GetPassword(), 6, 50); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := utils.ValidateUserRole(req.GetRole()); err != nil {
		violations = append(violations, fieldViolation("role", err))
	}
	return
}
