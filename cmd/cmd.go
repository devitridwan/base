package cmd

import (
	"base/cmd/consumer"
	"base/cmd/server"
	"base/internal/infrastructures/log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "Base Service",
		Short: "Base - Backend Service",
		Long:  "Base - API Gateway User Service",
	}
)

func Execute() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Register command in here
	rootCmd.AddCommand(server.ServeHTTPCmd())
	rootCmd.AddCommand(consumer.ServeConsumerCmd())
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln("Error: \n", err.Error())
		os.Exit(-1)
	}
}
