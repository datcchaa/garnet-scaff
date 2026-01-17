package server

import (
	"garnet-scaff/config"
	"garnet-scaff/internal/handler/api"
	"os"
	"os/signal"
	"syscall"

	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/log"
	"github.com/spf13/cobra"
)

var (
	serverHTTPCmd = &cobra.Command{
		Use:   "serve-http",
		Short: "garnet scaff Service HTTP",
		Long:  "garnet scaff Service HTTP",
		RunE:  run,
	}
)

func ServeHttpCmd() *cobra.Command {
	serverHTTPCmd.Flags().StringP("config", "c", "", "Config Path, both relative or absolute. i.e: /usr/local/bin/config/files")
	return serverHTTPCmd
}

func run(cmd *cobra.Command, args []string) error {
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

	server := api.New(&api.Options{
		Cfg:    appContainer.Cfg,
		UserUc: appContainer.UserUc,
	})

	go server.Run()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Info("Exiting gracefully...")
	case err := <-server.ListenError():
		log.Error("Error starting web server, exiting gracefully:", err)
	}

	return nil
}
