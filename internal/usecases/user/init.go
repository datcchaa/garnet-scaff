package user

import (
	"garnet-scaff/internal/domain/repository"
	"garnet-scaff/internal/usecases"

	"github.com/nocturna-ta/golib/txmanager"
)

type Module struct {
	userRepo repository.UserRepository
	txMgr    txmanager.TxManager
}
type Opts struct {
	UserRepo repository.UserRepository
	TxMgr    txmanager.TxManager
}

func New(opts *Opts) usecases.UserUseCases {
	return &Module{
		userRepo: opts.UserRepo,
		txMgr:    opts.TxMgr,
	}
}
