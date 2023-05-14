package main

import (
	"batchdispatcher/internal/dispatcher"
	"batchdispatcher/internal/job"
	"batchdispatcher/internal/logger"
	"batchdispatcher/internal/model"
	"batchdispatcher/internal/repository"
	"batchdispatcher/internal/server"
	"batchdispatcher/internal/timeutil"
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var l *zap.Logger
var configFilePath string // コンフィグファイルのパス

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceErrors: false,
	SilenceUsage:  false,
	RunE: func(cmd *cobra.Command, args []string) error {
		return start()
	},
}

func NewJob(name, batchCmd string) job.Job {
	return job.Job{
		Logger:           l,
		Name:             name,
		BatchCmd:         batchCmd,
		Status:           model.StatusNotRunning,
		LastChangeStatus: timeutil.NowFunc(),
	}
}

func NewDispatcher(jobs []*job.Job) (*dispatcher.Dispatcher, error) {
	d := &dispatcher.Dispatcher{
		Logger: l,
		Jobs:   jobs,
	}

	return d, nil
}

func NewServer(d *dispatcher.Dispatcher) (srv *server.Server, err error) {
	srv = &server.Server{}

	// set logger
	srv.Logger = l

	// set dispatcher
	srv.Dispatcher = d

	return srv, nil
}

func start() (err error) {
	l, err = logger.NewLogger()
	if err != nil {
		println(err)
		return err
	}
	l.Info("set logger")

	ctx := context.Background()
	configs, err := repository.LoadConfigFile(configFilePath)
	if err != nil {
		l.Error("failed to load config file", zap.Error(err))
		return err
	}

	l.Info("loaded config file")
	var targetJob []*job.Job
	for _, c := range configs {
		job := NewJob(c.Name, c.BatchCmd)
		targetJob = append(targetJob, &job)
	}

	d, err := NewDispatcher(targetJob)
	if err != nil {
		l.Error("failed to load dispatcher", zap.Error(err))
		return err
	}

	srv, err := NewServer(d)
	if err != nil {
		l.Error("failed to run server", zap.Error(err))
		return err
	}

	return srv.Start(ctx, "3000")
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVarP(&configFilePath, "config-file", "c", "./config.yaml", "config file path")
}
