package repository

import (
	"batchdispatcher/internal/model"
	"reflect"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantCs  []model.JobConfig
		wantErr bool
	}{
		{
			name: "load",
			args: args{
				path: "./conf_test.yaml",
			},
			wantCs: []model.JobConfig{
				{
					Name:     "ls",
					BatchCmd: "ls -la",
				},
				{
					Name:     "sleep",
					BatchCmd: "sleep 12s",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCs, err := LoadConfigFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCs, tt.wantCs) {
				t.Errorf("LoadConfigFile() = %v, want %v", gotCs, tt.wantCs)
			}
		})
	}
}
