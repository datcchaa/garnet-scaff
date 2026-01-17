package controller

import (
	_ "garnet-scaff/docs"
	"garnet-scaff/internal/usecases"
	"github.com/gofiber/swagger"
	"github.com/nocturna-ta/golib/router"
	"time"
)

type API struct {
	prefix         string
	port           uint
	readTimeout    time.Duration
	writeTimeout   time.Duration
	requestTimeout time.Duration
	enableSwagger  bool
	userUc         usecases.UserUseCases
}

type Options struct {
	Prefix         string
	Port           uint
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestTimeout time.Duration
	EnableSwagger  bool
	UserUc         usecases.UserUseCases
}

func New(opts *Options) *API {
	return &API{
		prefix:         opts.Prefix,
		port:           opts.Port,
		readTimeout:    opts.ReadTimeout,
		writeTimeout:   opts.WriteTimeout,
		requestTimeout: opts.RequestTimeout,
		enableSwagger:  opts.EnableSwagger,
		userUc:         opts.UserUc,
	}
}

func (api *API) RegisterRoute() *router.FastRouter {
	myRouter := router.New(&router.Options{
		Prefix:         api.prefix,
		Port:           api.port,
		ReadTimeout:    api.readTimeout,
		WriteTimeout:   api.writeTimeout,
		RequestTimeout: api.requestTimeout,
	})

	if api.enableSwagger {
		myRouter.CustomHandler("GET", "/docs/*", swagger.HandlerDefault, router.MustAuthorized(false))
	}

	myRouter.GET("/health", api.Ping, router.MustAuthorized(false))
	myRouter.Group("v1", func(v1 *router.FastRouter) {
		v1.Group("/users", func(user *router.FastRouter) {
			user.GET("", api.GetAllUsers, router.MustAuthorized(false))
			user.POST("", api.CreateUser, router.MustAuthorized(false))
			user.PATCH("", api.UpdateUser, router.MustAuthorized(false))
			user.GET("/:userId", api.GetUser, router.MustAuthorized(false))
			user.DELETE("/:userId", api.DeleteUser, router.MustAuthorized(false))
		})
	})

	return myRouter
}
