package config

import (
	"github.com/millisecond/adaptlb/model"
)

type Config struct {
	Frontends []*model.Frontend
	AWSConfig AWSConfig
}

type AWSConfig struct {
	Region          string
	Endpoint        string
	DynamoTableName string
	Key             string
	Secret          string
}
