package server

import (
	"garnet-scaff/config"
	"garnet-scaff/internal/handler/grpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/spf13/cobra"
)

var (
	serveGrpcCmd = &cobra.Command{
		Use:   "serve-grpc",
		Short: "garnet scaff Service gRPC",
		RunE:  runGrpc,
	}
)

func ServeGrpcCmd() *cobra.Command {
	serveGrpcCmd.Flags().StringP("config", "c", "", "Config Path, both relative or absolute. i.e: /usr/local/bin/config/files")
	return serveGrpcCmd
}

func runGrpc(cmd *cobra.Command, args []string) error {
	configLocation, _ := cmd.Flags().GetString("config")

	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, configLocation)

	database := sql.New(sql.DBConfig{
		SlaveDSN:        cfg.Database.SlaveDSN,
		MasterDSN:       cfg.Database.MasterDSN,
		RetryInterval:   cfg.Database.RetryInterval,
		MaxIdleConn:     cfg.Database.MaxIdleConn,
		MaxConn:         cfg.Database.MaxConn,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}, sql.DriverPostgres)

	appContainer := newContainer(&options{
		Cfg: cfg,
		DB:  database,
	})

	server := grpc.New(&grpc.Options{
		Cfg:    appContainer.Cfg,
		UserUc: appContainer.UserUc,
	})

	go server.Run()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Info("Exiting gracefully...")
	}

	return nil
}
