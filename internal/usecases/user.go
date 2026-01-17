package usecases

import (
	"context"
	"garnet-scaff/internal/usecases/request"
	"garnet-scaff/internal/usecases/response"

	"github.com/google/uuid"
)

type UserUseCases interface {
	GetAllUsers(ctx context.Context) (*[]response.UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error)
	RegisterUser(ctx context.Context, req *request.UserRegistrationRequest) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, req *request.UserUpdateRequest) (*response.UserResponse, error)
	SoftDeleteUser(ctx context.Context, id uuid.UUID) error
}
