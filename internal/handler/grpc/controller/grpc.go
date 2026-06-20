package controller

import (
	"garnet-scaff/internal/handler/grpc/proto"
	"garnet-scaff/internal/usecases"

	"github.com/nocturna-ta/golib/grpc"
)

type GRPC struct {
	port   uint
	userUc usecases.UserUseCases
}

type server struct {
	proto.UnimplementedUserServiceServer
	userUc usecases.UserUseCases
}

type Options struct {
	Port   uint
	UserUc usecases.UserUseCases
}

func New(opts *Options) *GRPC {
	return &GRPC{
		port:   opts.Port,
		userUc: opts.UserUc,
	}
}

func (g *GRPC) Register() *grpc.Server {
	srv := grpc.NewServer(&grpc.ServerOptions{
		Port: g.port,
	})

	grpcServiceServer := &server{
		userUc: g.userUc,
	}

	srv.Register(proto.RegisterUserServiceServer, grpcServiceServer)

	return srv
}
