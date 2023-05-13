package job

import (
	"batchdispatcher/internal/model"
	"batchdispatcher/internal/timeutil"
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

var l *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, _ = config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
}
func TestJob_Run(t *testing.T) {
	type fields struct {
		Logger           *zap.Logger
		Name             string
		BatchCmd         string
		Status           model.JobStatus
		LastChangeStatus time.Time
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ls",
			fields: fields{
				Logger:           l,
				Name:             "ls",
				BatchCmd:         "ls",
				Status:           model.StatusComplete,
				LastChangeStatus: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "xxx (not found)",
			fields: fields{
				Logger:           l,
				Name:             "xxx",
				BatchCmd:         "xxx",
				Status:           model.StatusComplete,
				LastChangeStatus: time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set test time
			timeutil.NowFunc = func() time.Time {
				return time.Date(2000, 1, 23, 0, 0, 0, 0, time.Local)
			}

			j := &Job{
				Logger:           tt.fields.Logger,
				Name:             tt.fields.Name,
				BatchCmd:         tt.fields.BatchCmd,
				Status:           tt.fields.Status,
				LastChangeStatus: tt.fields.LastChangeStatus,
			}
			if err := j.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Job.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
