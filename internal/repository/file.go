package repository

import (
	"batchdispatcher/internal/model"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfigFile(path string) (cs []model.JobConfig, err error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return
	}

	data, err := readOnStruct(buf)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readOnStruct(fileBuffer []byte) ([]model.JobConfig, error) {
	data := make([]model.JobConfig, 0)
	err := yaml.Unmarshal(fileBuffer, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
