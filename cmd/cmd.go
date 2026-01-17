package cmd

import (
	"garnet-scaff/cmd/server"
	"os"

	"github.com/nocturna-ta/golib/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Garnet Scaff",
		Short: "Garnet scaff",
	}
)

func Execute() {
	log.SetFormatter("json")
	rootCmd.AddCommand(server.ServeHttpCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error: ", err.Error())
		os.Exit(-1)
	}
}
