package job

import (
	"batchdispatcher/internal/model"
	"batchdispatcher/internal/timeutil"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Job struct {
	Logger           *zap.Logger
	Name             string
	BatchCmd         string
	Status           model.JobStatus
	LastChangeStatus time.Time
	LastSucessStatus time.Time
}

func (j *Job) Run(ctx context.Context) error {
	j.ChangeStatus(model.StatusRunning)
	cmdName, cmdArgs := separateNameArgs(j.BatchCmd)
	cmd := exec.Command(cmdName, cmdArgs...)
	err := cmd.Run()

	if err != nil {
		j.ChangeStatus(model.StatusError)
		j.Logger.Error("exec error", zap.Error(err))
		return err
	}

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		j.ChangeStatus(model.StatusError)
		j.Logger.Error("unexpected exit code", zap.Error(err))
		return fmt.Errorf("unexpected exit code")
	}

	j.ChangeStatus(model.StatusComplete)
	return nil
}

func (j *Job) ChangeStatus(nextStatus model.JobStatus) {
	j.Status = nextStatus
	t := timeutil.NowFunc()
	if j.Status == model.StatusComplete {
		j.LastSucessStatus = t
	}
	j.LastChangeStatus = t
	j.Logger.Info("change status", zap.String("name", j.Name), zap.String("status", string(j.Status)))
}

func (j *Job) GetName() (jobName string) {
	return j.Name
}

func (j *Job) GetBatchCmd() (batchCmd string) {
	return j.BatchCmd
}

func (j *Job) GetStatus() (status model.JobStatus) {
	return j.Status
}

func (j *Job) GetLastChangeStatus() time.Time {
	return j.LastChangeStatus
}

func (j *Job) GetLastSucessStatus() time.Time {
	return j.LastSucessStatus
}

func separateNameArgs(cmdStr string) (cmdName string, cmdArgs []string) {
	sep := strings.Split(cmdStr, " ")
	if len(sep) <= 1 {
		return sep[0], []string{}
	}
	return sep[0], sep[1:]
}
