package controller

import (
	"context"
	"garnet-scaff/internal/handler/grpc/proto"
	"garnet-scaff/internal/usecases/request"

	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/tracing"
)

func (s *server) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetAllUser")
	defer span.End()

	users, err := s.userUc.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var protoUsers []*proto.UserResponse
	for _, user := range users {
		protoUsers = append(protoUsers, &proto.UserResponse{
			Id:       user.ID,
			Username: user.Username,
		})
	}

	return &proto.GetAllUsersResponse{
		Users: protoUsers,
	}, nil
}

func (s *server) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.UserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.GetUserByID")
	defer span.End()

	userId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	user, err := s.userUc.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *server) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.UserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.RegisterUser")
	defer span.End()

	userReq := &request.UserRegistrationRequest{
		Username: req.Username,
	}

	user, err := s.userUc.RegisterUser(ctx, userReq)
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.UpdateUser")
	defer span.End()

	userReq := &request.UserUpdateRequest{
		ID:       req.Id,
		Username: req.Username,
	}

	user, err := s.userUc.UpdateUser(ctx, userReq)
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (s *server) SoftDeleteUser(ctx context.Context, req *proto.SoftDeleteUserRequest) (*proto.SoftDeleteUserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "Controller.SoftDeleteUser")
	defer span.End()

	userId, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	err = s.userUc.SoftDeleteUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &proto.SoftDeleteUserResponse{}, nil
}
