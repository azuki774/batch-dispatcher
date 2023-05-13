package dispatcher

import (
	"batchdispatcher/internal/job"
	"context"
	"testing"

	"go.uber.org/zap"
)

var l *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, _ = config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
}

func TestDispatcher_Run(t *testing.T) {
	type fields struct {
		Logger *zap.Logger
		Jobs   []*job.Job
	}
	type args struct {
		ctx     context.Context
		jobname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Logger: l,
				Jobs: []*job.Job{
					{
						Logger: l,
						Name:   "ls",
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				jobname: "ls",
			},
			wantErr: false,
		},
		{
			name: "not found",
			fields: fields{
				Logger: l,
				Jobs: []*job.Job{
					{
						Name: "xxx",
					},
				},
			},
			args: args{
				ctx:     context.Background(),
				jobname: "not",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{
				Logger: tt.fields.Logger,
				Jobs:   tt.fields.Jobs,
			}
			if err := d.Run(tt.args.ctx, tt.args.jobname); (err != nil) != tt.wantErr {
				t.Errorf("Dispatcher.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
