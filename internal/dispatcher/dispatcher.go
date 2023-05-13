package dispatcher

import (
	"batchdispatcher/internal/job"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Dispatcher struct {
	Logger *zap.Logger
	Jobs   []job.Job
}

func (d *Dispatcher) Run(ctx context.Context, jobname string) (err error) {
	for _, j := range d.Jobs {
		if j.Name == jobname {
			// Job nameが一致する場合のみ実行
			if j.GetStatus() == job.StatusRunning {
				// already running
				return fmt.Errorf("already running")
			}

			// job Run
			go j.Run(ctx)
			return nil
		}
	}

	// job not found
	return fmt.Errorf("job not found")
}
