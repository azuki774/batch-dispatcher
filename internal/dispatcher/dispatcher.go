package dispatcher

import "batchdispatcher/internal/job"

type Dispatcher struct {
	Jobs []job.Job
}
