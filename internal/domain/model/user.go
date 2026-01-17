package model

import (
	"garnet-scaff/internal/usecases/request"
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseModel
	ID       uuid.UUID `db:"id" json:"id"`
	Username string    `db:"username" json:"username"`
}

func ConstructRegistrationUser(req *request.UserRegistrationRequest) *User {
	now := time.Now()
	user := &User{
		BaseModel: BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
			IsDeleted: false,
		},
		ID:       uuid.New(),
		Username: req.Username,
	}

	return user
}
