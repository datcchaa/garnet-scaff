package dao

import (
	"context"
	sql2 "database/sql"
	"errors"
	"garnet-scaff/internal/domain/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/txmanager/utils"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Insert(t *testing.T) {
	db, mockDb, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dataStore := &sql.Store{
		Master: &sql.DB{DBConnection: sqlxDB},
		Slave:  &sql.DB{DBConnection: sqlxDB},
	}

	mockDb.ExpectBegin()

	sqlxTx, err := sqlxDB.BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	id := uuid.New()
	username := "test"
	now := time.Now()

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		fn      func()
	}{
		{
			name: "ShouldError_Duplicate",
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID:       id,
					Username: username,
					BaseModel: model.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
						IsDeleted: false,
					},
				},
			},
			wantErr: true,
			fn: func() {
				mockDb.ExpectExec(`INSERT INTO users`).WillReturnError(&pq.Error{Code: "23505"})
			},
		},
		{
			name: "ShouldError_FailedInsert",
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID:       id,
					Username: username,
				},
			},
			wantErr: true,
			fn: func() {
				mockDb.ExpectExec(`INSERT INTO users`).WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "Success",
			args: args{
				ctx: utils.SetSqlTx(context.Background(), sqlxTx),
				user: &model.User{
					ID:       id,
					Username: username,
				},
			},
			wantErr: false,
			fn: func() {
				mockDb.ExpectExec(`INSERT INTO users`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			a := NewUserRepository(&OptsUserRepository{
				DB: dataStore,
			})

			if err := a.Insert(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mockDb, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dataStore := &sql.Store{
		Master: &sql.DB{DBConnection: sqlxDB},
		Slave:  &sql.DB{DBConnection: sqlxDB},
	}

	mockDb.ExpectBegin()

	_, err := sqlxDB.BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	id := uuid.New()
	username := "test"

	excpected := &model.User{
		ID:       id,
		Username: username,
	}

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		fn      func()
	}{
		{
			name: "ShouldError_WhenFailedUpdate",
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: excpected.ID,
				},
			},
			wantErr: true,
			fn: func() {
				mockDb.ExpectExec(`UPDATE users SET`).WillReturnError(errors.New("failed update"))
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: excpected.ID,
				},
			},
			wantErr: false,
			fn: func() {
				mockDb.ExpectExec(`UPDATE users SET`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			a := NewUserRepository(&OptsUserRepository{
				DB: dataStore,
			})
			if err := a.Update(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, mockDb, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	dataStore := &sql.Store{
		Master: &sql.DB{DBConnection: sqlxDB},
		Slave:  &sql.DB{DBConnection: sqlxDB},
	}

	mockDb.ExpectBegin()

	_, err := sqlxDB.BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	id := uuid.New()
	excpected := &model.User{
		ID: id,
	}

	type args struct {
		ctx  context.Context
		user *model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		fn      func()
	}{
		{
			name: "ShouldError_FailedSoftDelete",
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: excpected.ID,
				},
			},
			wantErr: true,
			fn: func() {
				mockDb.ExpectExec(`UPDATE users SET`).WillReturnError(errors.New("failed soft delete"))
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				user: &model.User{
					ID: excpected.ID,
				},
			},
			wantErr: false,
			fn: func() {
				mockDb.ExpectExec(`UPDATE users SET`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			a := NewUserRepository(&OptsUserRepository{
				DB: dataStore,
			})
			if err := a.Delete(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mockDb, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	dataStore := &sql.Store{
		Master: &sql.DB{DBConnection: sqlxDB},
		Slave:  &sql.DB{DBConnection: sqlxDB},
	}

	mockDb.ExpectBegin()

	_, err := sqlxDB.BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	testName := "test"
	expected := model.User{
		ID:       uuid.New(),
		Username: testName,
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		fn      func()
	}{
		{
			name: "ShouldError_WhenFailedFind",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			wantErr: true,
			fn: func() {
				mockDb.ExpectQuery("SELECT users.id, users.username, users.created_at, users.updated_at FROM users").
					WillReturnError(errors.New("failed find id"))
			},
		},
		{
			name: "ShouldError_WhenDocumentExists",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			wantErr: true,
			fn: func() {
				mockDb.ExpectQuery("SELECT users.id, users.username, users.created_at, users.updated_at FROM users").
					WillReturnError(sql2.ErrNoRows)
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			wantErr: false,
			fn: func() {
				rows := mockDb.NewRows([]string{
					"id", "username",
				}).AddRow(expected.ID, expected.Username)
				mockDb.ExpectQuery("SELECT users.id, users.username, users.created_at, users.updated_at FROM users").WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			a := NewUserRepository(&OptsUserRepository{
				DB: dataStore,
			})
			_, err := a.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	db, mockDb, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	dataStore := &sql.Store{
		Master: &sql.DB{DBConnection: sqlxDB},
		Slave:  &sql.DB{DBConnection: sqlxDB},
	}

	mockDb.ExpectBegin()

	_, err := sqlxDB.BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	testName := "test"
	expected := []*model.User{
		{
			ID:       uuid.New(),
			Username: testName,
		},
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		fn      func()
	}{
		{
			name: "ShouldError_WhenFailedFind",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			fn: func() {
				mockDb.ExpectQuery("SELECT users.id, users.username, users.created_at, users.updated_at FROM users").WillReturnError(errors.New("failed find all"))
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			fn: func() {
				rows := mockDb.NewRows([]string{
					"id", "username",
				}).AddRow(expected[0].ID, expected[0].Username)
				mockDb.ExpectQuery("SELECT users.id, users.username, users.created_at, users.updated_at FROM users").WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			a := NewUserRepository(&OptsUserRepository{
				DB: dataStore,
			})
			_, err := a.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
