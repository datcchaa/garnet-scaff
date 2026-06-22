package user

import (
	"context"
	"errors"
	"garnet-scaff/internal/domain/model"
	"garnet-scaff/internal/domain/repository/mocks_repository"
	"garnet-scaff/internal/usecases/request"
	"garnet-scaff/internal/usecases/response"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestModule_GetUserByID(t *testing.T) {
	mockedUserRepo := &mocks_repository.UserRepository{}
	expectedUser := &model.User{
		ID: uuid.New(),
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name     string
		args     args
		fn       func()
		wantResp *response.UserResponse
		wantErr  bool
	}{
		{
			name: "ShouldError_WhenFailedGetUser",
			args: args{
				ctx: context.Background(),
				id:  expectedUser.ID,
			},
			fn: func() {
				mockedUserRepo.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("failed")).Once()
			},
			wantResp: nil,
			wantErr:  true,
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  expectedUser.ID,
			},
			fn: func() {
				mockedUserRepo.On("FindByID", mock.Anything, mock.Anything).Return(expectedUser, nil).Once()
			},
			wantResp: &response.UserResponse{
				ID:       expectedUser.ID.String(),
				Username: expectedUser.Username,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()

			m := New(&Opts{
				UserRepo: mockedUserRepo,
			})
			_, err := m.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestModule_RegisterUser(t *testing.T) {
	mockUserRepo := &mocks_repository.UserRepository{}
	expectedUser := &model.User{
		ID: uuid.New(),
	}

	db
	type args struct {
		ctx context.Context
		req *request.UserRegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *response.UserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				userRepo: tt.fields.userRepo,
				txMgr:    tt.fields.txMgr,
			}
			got, err := m.RegisterUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
