package server

import (
	"context"
	"garnet-scaff/config"
	"garnet-scaff/internal/interfaces/dao"
	"garnet-scaff/internal/usecases"
	"garnet-scaff/internal/usecases/user"

	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/txmanager"

	txSql "github.com/nocturna-ta/golib/txmanager/sql"
)

type container struct {
	Cfg    config.MainConfig
	UserUc usecases.UserUseCases
}

type options struct {
	Cfg *config.MainConfig
	DB  *sql.Store
}

func newContainer(opts *options) *container {
	usersRepo := dao.NewUserRepository(&dao.OptsUserRepository{
		DB: opts.DB,
	})

	txMgr, err := txmanager.New(context.Background(), &txmanager.DriverConfig{
		Type: "sql",
		Config: txSql.Config{
			DB: opts.DB,
		},
	})
	if err != nil {
		log.Fatal("Failed to instantiate transaction manager ")
	}

	userUc := user.New(&user.Opts{
		UserRepo: usersRepo,
		TxMgr:    txMgr,
	})

	return &container{
		Cfg:    *opts.Cfg,
		UserUc: userUc,
	}

}
