package repository

import (
	"context"
	"garnet-scaff/internal/domain/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}
