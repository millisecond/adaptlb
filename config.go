package adaptlb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type AWSConfig struct {
	Region          string
	Endpoint        string
	DynamoTableName string
	Key             string
	Secret          string
}

func (cfg *AWSConfig) Generate() *aws.Config {
	c := &aws.Config{}
	if len(cfg.Region) > 0 {
		c.Region = aws.String(cfg.Region)
	}
	if len(cfg.Endpoint) > 0 {
		c.Endpoint = aws.String(cfg.Endpoint)
	}
	// Only set credentials if key and secret are given, otherwise fall back to IAM role
	if len(cfg.Key) > 0 && len(cfg.Secret) > 0 {
		c.Credentials = credentials.NewStaticCredentials(cfg.Key, cfg.Secret, "")
	}
	return c
}
