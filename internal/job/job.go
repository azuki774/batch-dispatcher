package job

import (
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
	Status           JobStatus
	LastChangeStatus time.Time
}

func NewJob(name, batchDir string) Job {
	return Job{
		Name:             name,
		BatchDir:         batchDir,
		Status:           StatusNotRunning,
		LastChangeStatus: time.Now(),
	}
}
