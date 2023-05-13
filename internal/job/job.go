package job

import (
	"batchdispatcher/internal/timeutil"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type JobStatus string

var (
	StatusNotRunning = JobStatus("NOTRUNNING")
	StatusRunning    = JobStatus("RUNNING")
	StatusError      = JobStatus("ERROR")
	StatusComplete   = JobStatus("COMPLETE")
)

type Job struct {
	Name             string
	BatchCmd         string
	status           JobStatus
	LastChangeStatus time.Time
}

func NewJob(name, batchCmd string) Job {
	return Job{
		Name:             name,
		BatchCmd:         batchCmd,
		status:           StatusNotRunning,
		LastChangeStatus: timeutil.NowFunc(),
	}
}

func (j *Job) Run(ctx context.Context) error {
	j.ChangeStatus(StatusRunning)
	cmd := exec.Command(j.BatchCmd)
	err := cmd.Run()

	if err != nil {
		j.ChangeStatus(StatusError)
		return err
	}

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		j.ChangeStatus(StatusError)
		return fmt.Errorf("unexpected exit code")

	}

	j.ChangeStatus(StatusComplete)
	return nil
}

func (j *Job) ChangeStatus(nextStatus JobStatus) {
	j.status = nextStatus
	j.LastChangeStatus = timeutil.NowFunc()
}

func (j *Job) GetStatus() (status JobStatus) {
	return j.status
}
