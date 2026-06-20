package user

import (
	"context"
	"errors"
	"garnet-scaff/internal/domain/model"
	"garnet-scaff/internal/interfaces/dao"
	"garnet-scaff/internal/usecases/request"
	"garnet-scaff/internal/usecases/response"

	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/log"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
)

func (m *Module) GetUserByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserUseCases.GetUserByID")
	defer span.End()

	user, err := m.userRepo.FindByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[UserUseCases.GetUserByID] failed to find user")
		return nil, &custerr.ErrChain{
			Message: "failed to find user",
			Code:    500,
			Type:    response2.ErrInternalServerError,
			Cause:   err,
		}
	}

	return &response.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}

func (m *Module) RegisterUser(ctx context.Context, req *request.UserRegistrationRequest) (*response.UserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserUseCases.RegisterUser")
	defer span.End()

	var (
		user *model.User
	)

	user = model.ConstructRegistrationUser(req)

	err := m.userRepo.Insert(ctx, user)
	if err != nil {
		if errors.Is(err, dao.ErrDuplicate) {
			return nil, &custerr.ErrChain{
				Message: "user is already registered",
				Code:    400,
				Type:    response2.ErrBadRequest,
				Cause:   err,
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"user":  user,
		}).ErrorWithCtx(ctx, "[UserUseCases.RegisterUser] failed to insert user")
		return nil, err
	}

	return &response.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}

func (m *Module) UpdateUser(ctx context.Context, req *request.UserUpdateRequest) (*response.UserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserUseCases.UpdateUser")
	defer span.End()

	var (
		existing *model.User
	)

	id, err := uuid.Parse(req.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    req.ID,
		}).ErrorWithCtx(ctx, "[UserUseCases.UpdateUser] failed to parse id")
		return nil, &custerr.ErrChain{
			Message: "failed to parse id",
			Code:    500,
			Type:    response2.ErrInternalServerError,
			Cause:   err,
		}
	}

	transaction := func(txCtx context.Context) (any, error) {
		var errTx error
		existing, errTx = m.userRepo.FindByID(ctx, id)
		if errTx != nil {
			log.WithFields(log.Fields{
				"error": errTx,
				"id":    id,
			}).ErrorWithCtx(txCtx, "[UserUseCases.UpdateUser] failed to find user")
			return nil, &custerr.ErrChain{
				Message: "failed to find user",
				Code:    500,
				Type:    response2.ErrInternalServerError,
				Cause:   errTx,
			}
		}

		existing.Username = req.Username
		errTx = m.userRepo.Update(ctx, existing)
		if errTx != nil {
			if errors.Is(errTx, dao.ErrNoUpdateHappened) {
				return nil, &custerr.ErrChain{
					Message: "failed to update user",
					Code:    500,
					Type:    response2.ErrInternalServerError,
					Cause:   errTx,
				}
			}
			log.WithFields(log.Fields{
				"error": errTx,
				"user":  existing,
			}).ErrorWithCtx(txCtx, "[UserUseCases.UpdateUser] failed to update user")
			return nil, &custerr.ErrChain{
				Message: "failed to update user",
				Code:    500,
				Type:    response2.ErrInternalServerError,
				Cause:   errTx,
			}
		}
		return nil, nil
	}

	_, err = m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"user":  existing,
		}).ErrorWithCtx(ctx, "[UserUseCases.UpdateUser] failed to execute transaction")
		return nil, err
	}

	return &response.UserResponse{
		ID:       existing.ID.String(),
		Username: existing.Username,
	}, nil
}

func (m *Module) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserUseCases.SoftDeleteUser")
	defer span.End()

	var (
		existing *model.User
	)

	transaction := func(txCtx context.Context) (any, error) {
		var errTx error
		existing, errTx = m.userRepo.FindByID(ctx, id)
		if errTx != nil {
			log.WithFields(log.Fields{
				"error": errTx,
				"id":    id,
			}).ErrorWithCtx(txCtx, "[UserUseCases.SoftDeleteUser] failed to find user")
			return nil, &custerr.ErrChain{
				Message: "failed to find user",
				Code:    500,
				Type:    response2.ErrInternalServerError,
				Cause:   errTx,
			}
		}

		existing.Username = ""

		errTx = m.userRepo.Delete(ctx, existing)
		if errTx != nil {
			if errors.Is(errTx, dao.ErrNoUpdateHappened) {
				return nil, &custerr.ErrChain{
					Message: "failed to delete user",
					Code:    500,
					Type:    response2.ErrInternalServerError,
					Cause:   errTx,
				}
			}
			log.WithFields(log.Fields{
				"error": errTx,
				"user":  existing,
			}).ErrorWithCtx(txCtx, "[UserUseCases.SoftDeleteUser] failed to delete user")
			return nil, &custerr.ErrChain{
				Message: "failed to delete user",
				Code:    500,
				Type:    response2.ErrInternalServerError,
				Cause:   errTx,
			}
		}

		return nil, nil
	}

	_, err := m.txMgr.Execute(ctx, transaction, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"user":  existing,
		}).ErrorWithCtx(ctx, "[UserUseCases.SoftDeleteUser] failed to execute transaction")
		return err
	}

	return nil
}

func (m *Module) GetAllUsers(ctx context.Context) ([]*response.UserResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserUseCases.GetAllUsers")
	defer span.End()

	users, err := m.userRepo.GetAll(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[UserUseCases.GetAllUsers] failed to get all users")
		return nil, &custerr.ErrChain{
			Message: "failed to get all users",
			Code:    500,
			Type:    response2.ErrInternalServerError,
			Cause:   err,
		}
	}

	var userResponse []*response.UserResponse
	for _, user := range users {
		userResp := response.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
		}

		userResponse = append(userResponse, &userResp)
	}

	return userResponse, nil
}
