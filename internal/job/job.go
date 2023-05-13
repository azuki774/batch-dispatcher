package job

import (
	"batchdispatcher/internal/timeutil"
	"context"
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
	BatchDir         string
	status           JobStatus
	LastChangeStatus time.Time
}

func NewJob(name, batchDir string) Job {
	return Job{
		Name:             name,
		BatchDir:         batchDir,
		status:           StatusNotRunning,
		LastChangeStatus: timeutil.NowFunc(),
	}
}

func (j *Job) Run(ctx context.Context) {
	j.ChangeStatus(StatusRunning)
	cmd := exec.Command(j.BatchDir)
	err := cmd.Run()

	if err != nil {
		j.ChangeStatus(StatusError)
	}

	exitCode := cmd.ProcessState.ExitCode()
	if exitCode != 0 {
		j.ChangeStatus(StatusError)

	}

	j.ChangeStatus(StatusComplete)
}

func (j *Job) ChangeStatus(nextStatus JobStatus) {
	j.status = nextStatus
	j.LastChangeStatus = timeutil.NowFunc()
}

func (j *Job) GetStatus() (status JobStatus) {
	return j.status
}
