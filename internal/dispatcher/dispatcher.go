package dispatcher

import (
	"batchdispatcher/internal/job"
	"batchdispatcher/internal/model"
	"context"
	"fmt"

	"go.uber.org/zap"
)

// type Job interface {
// 	Run(ctx context.Context) error
// 	ChangeStatus(nextStatus model.JobStatus)

// 	GetName() (jobName string)
// 	GetBatchCmd() (batchCmd string)
// 	GetStatus() (status model.JobStatus)
// 	GetLastChangeStatus() time.Time
// 	GetLastSucessStatus() time.Time
// }

type Dispatcher struct {
	Logger *zap.Logger
	Jobs   []*job.Job
}

func (d *Dispatcher) Run(ctx context.Context, jobname string) (err error) {
	for _, j := range d.Jobs {
		if j.GetName() == jobname {
			// Job nameが一致する場合のみ実行
			if j.GetStatus() == model.StatusRunning {
				// already running
				return model.ErrAlreadyRunning
			}

			// job Run
			go j.Run(ctx)
			return nil
		}
	}

	// job not found
	return model.ErrNotFound
}

// GetJobsInfo は Job処理用の構造体から表示用の構造体に変換
func (d *Dispatcher) GetJobsInfo() (jobInfo []model.JobInfo) {
	for _, job := range d.Jobs {
		fmt.Println(job)
		jobInfo = append(jobInfo, model.JobInfo{
			Name:             job.GetName(),
			BatchCmd:         job.GetBatchCmd(),
			Status:           job.GetStatus(),
			LastChangeStatus: job.GetLastChangeStatus(),
			LastSucessStatus: job.GetLastSucessStatus(),
		})
	}

	return jobInfo
}
