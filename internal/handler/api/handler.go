package api

import (
	"garnet-scaff/config"
	"garnet-scaff/internal/usecases"

	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/router"
)

type Options struct {
	Cfg    config.MainConfig
	UserUc usecases.UserUseCases
}

type Handler struct {
	opts        *Options
	listerErrCh chan error
	myRouter    *router.FastRouter
}

func New(opts *Options) *Handler {
	handler := &Handler{
		opts: opts,
	}

	handler.myRouter = nil

	return handler
}

func (h *Handler) Run() {
	log.Info("API Listening on %d", h.opts.Cfg.Server.Port)
	h.listerErrCh <- h.myRouter.StartServe()
}

func (h *Handler) ListenError() <-chan error {
	return h.listerErrCh
}
