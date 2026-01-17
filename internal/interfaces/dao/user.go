package dao

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	"garnet-scaff/internal/domain/model"
	"garnet-scaff/internal/domain/repository"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/golib/txmanager/utils"
)

type UserRepository struct {
	db *sql.Store
}

type OptsUserRepository struct {
	DB *sql.Store
}

func NewUserRepository(opts *OptsUserRepository) repository.UserRepository {
	return &UserRepository{
		db: opts.DB,
	}
}

const (
	insertUser = `INSERT INTO users (id, username, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	selectUser = `SELECT %s FROM users %s WHERE TRUE %s`
	updateUser = `UPDATE users SET %s WHERE TRUE %s`
)

func (repo *UserRepository) Insert(ctx context.Context, user *model.User) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserRepository.Insert")
	defer span.End()

	var (
		err error
	)

	sqlTrx := utils.GetSqlTx(ctx)
	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, insertUser,
			user.ID,
			user.Username,
			user.CreatedAt,
			user.UpdatedAt)
	} else {
		_, err = repo.db.GetMaster().ExecContext(ctx, insertUser,
			user.ID,
			user.Username,
			user.CreatedAt,
			user.UpdatedAt)
	}

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				log.WithFields(log.Fields{
					"error": err,
					"user":  user,
				}).ErrorWithCtx(ctx, "[UserRepository.Insert] User already exists")
				return ErrDuplicate
			}
		}
		log.WithFields(log.Fields{
			"error": err,
			"user":  user,
		}).ErrorWithCtx(ctx, "[UserRepository.Insert] Failed to insert user")
		return err
	}
	return nil
}

func (repo *UserRepository) Update(ctx context.Context, user *model.User) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserRepository.Update")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err  error
		args []any
	)

	setQuery := "username = $1, updated_at = $2"
	whereQuery := " AND id = $3 AND is_deleted = FALSE"

	args = append(args, user.Username, time.Now(), user.ID)
	query := fmt.Sprintf(updateUser, setQuery, whereQuery)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		_, err = repo.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"user":  user,
		}).ErrorWithCtx(ctx, "[UserRepository.Update] Failed to update user")
		return err
	}

	return nil
}

func (u *UserRepository) Delete(ctx context.Context, user *model.User) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserRepository.Delete")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		err  error
		args []any
	)

	setQuery := "is_deleted = TRUE"
	whereQuery := " AND id = $1 AND is_deleted = FALSE"

	args = append(args, user.ID)
	query := fmt.Sprintf(updateUser, setQuery, whereQuery)

	if sqlTrx != nil {
		_, err = sqlTrx.ExecContext(ctx, query, args...)
	} else {
		_, err = u.db.GetMaster().ExecContext(ctx, query, args...)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"user":  user,
		}).ErrorWithCtx(ctx, "[UserRepository.Delete] Failed to delete user")
		return err
	}

	return nil
}

func (u *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserRepository.FindByID")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)

	var (
		user model.User
		err  error
		args []any
	)

	selectQuery := "users.id, users.username, users.created_at, users.updated_at"
	whereQuery := " AND users.id = $1 AND is_deleted = FALSE"

	args = append(args, id)

	query := fmt.Sprintf(selectUser, selectQuery, "", whereQuery)

	if sqlTrx != nil {
		err = sqlTrx.GetContext(ctx, &user, query, args...)
	} else {
		err = u.db.GetMaster().GetContext(ctx, &user, query, args...)
	}

	if err != nil {
		if errors.Is(err, sql2.ErrNoRows) {
			return nil, ErrNoResult
		}
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[UserRepository.FindByID] Failed to find user")
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "UserRepository.GetAll")
	defer span.End()

	sqlTrx := utils.GetSqlTx(ctx)
	var (
		users []*model.User
		err   error
	)

	selectQuery := "users.id, users.username, users.created_at, users.updated_at"
	whereQuery := " AND users.is_deleted = FALSE"

	query := fmt.Sprintf(selectUser, selectQuery, "", whereQuery)
	if sqlTrx != nil {
		err = sqlTrx.SelectContext(ctx, &users, query)
	} else {
		err = u.db.GetMaster().SelectContext(ctx, &users, query)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[UserRepository.GetAll] Failed to get users")
		return nil, err
	}

	return users, err
}
