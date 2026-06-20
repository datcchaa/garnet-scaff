package grpc

import (
	"garnet-scaff/config"
	"garnet-scaff/internal/handler/grpc/controller"
	"garnet-scaff/internal/usecases"

	"github.com/nocturna-ta/golib/grpc"
)

type Options struct {
	Cfg    config.MainConfig
	UserUc usecases.UserUseCases
}

type Handler struct {
	opts       *Options
	grpcServer *grpc.Server
}

func New(opts *Options) *Handler {
	handler := &Handler{
		opts: opts,
	}

	handler.grpcServer = controller.New(&controller.Options{
		Port:   opts.Cfg.GrpcServer.Port,
		UserUc: opts.UserUc,
	}).Register()

	return handler
}

func (h *Handler) Run() {
	h.grpcServer.MustStart()
}

func (h *Handler) Stop() {
	h.grpcServer.Stop()
}
