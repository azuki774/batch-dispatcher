package main

import (
	"batchdispatcher/internal/dispatcher"
	"batchdispatcher/internal/job"
	"batchdispatcher/internal/logger"
	"batchdispatcher/internal/model"
	"batchdispatcher/internal/server"
	"batchdispatcher/internal/timeutil"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

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

func NewJob(name, batchCmd string) job.Job {
	l, _ := logger.NewLogger()
	return job.Job{
		Logger:           l,
		Name:             name,
		BatchCmd:         batchCmd,
		Status:           model.StatusNotRunning,
		LastChangeStatus: timeutil.NowFunc(),
	}
}

func NewDispatcher(jobs []*job.Job) (*dispatcher.Dispatcher, error) {
	l, err := logger.NewLogger()
	if err != nil {
		return &dispatcher.Dispatcher{}, err
	}

	d := &dispatcher.Dispatcher{
		Logger: l,
		Jobs:   jobs,
	}

	return d, nil
}

func NewServer(d *dispatcher.Dispatcher) (srv *server.Server, err error) {
	srv = &server.Server{}

	// set logger
	l, err := logger.NewLogger()
	if err != nil {
		return &server.Server{}, err
	}
	srv.Logger = l

	// set dispatcher
	l.Info("set dispatcher")
	srv.Dispatcher = d

	return srv, nil
}

func start() (err error) {
	ctx := context.Background()
	samplejobs1 := NewJob("ls", "ls -la")
	samplejobs2 := NewJob("sleep", "sleep 15s")
	d, err := NewDispatcher([]*job.Job{&samplejobs1, &samplejobs2})
	if err != nil {
		fmt.Println(err)
		return err
	}

	srv, err := NewServer(d)
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
