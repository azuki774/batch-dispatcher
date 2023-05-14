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
	Name              string    `json:"name"`
	BatchCmd          string    `json:"cmd"`
	Status            JobStatus `json:"status"`
	LastChangeStatus  time.Time `json:"last_change_status"`
	LastSuccessStatus time.Time `json:"last_success_status"`
}

type JobConfig struct {
	Name     string `yaml:"name"`
	BatchCmd string `yaml:"cmd"`
}
