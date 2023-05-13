package model

import (
	"errors"
	"time"
)

type JobStatus string

var (
	StatusNotRunning = JobStatus("NOTRUNNING")
	StatusRunning    = JobStatus("RUNNING")
	StatusError      = JobStatus("ERROR")
	StatusComplete   = JobStatus("COMPLETE")
)

var ErrNotFound = errors.New("job not found")
var ErrAlreadyRunning = errors.New("already running")

type JobInfo struct {
	Name             string
	BatchCmd         string
	Status           JobStatus
	LastChangeStatus time.Time
	LastSucessStatus time.Time
}
