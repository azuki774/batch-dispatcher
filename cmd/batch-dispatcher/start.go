package main

import (
	"batchdispatcher/internal/server"
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type StartOption struct {
	Logger *zap.Logger
}

var (
	version  string
	revision string
	build    string
)

var startOpt StartOption

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return start()
	},
}

func start() (err error) {
	ctx := context.Background()
	srv, err := server.NewServer()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return srv.Start(ctx, "3000")
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
