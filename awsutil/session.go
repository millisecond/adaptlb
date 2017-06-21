package awsutil

import (
	"github.com/millisecond/adaptlb/config"

	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/containous/traefik/log"
	"github.com/facebookgo/ensure"
	"testing"
)

func CreateSession(cfg *config.Config) *session.Session {
	return session.Must(session.NewSession(CreateConfig(cfg)))
}

func CreateConfig(cfg *config.Config) *aws.Config {
	c := &aws.Config{}
	if len(cfg.AWSConfig.Region) > 0 {
		c.Region = aws.String(cfg.AWSConfig.Region)
	}
	if len(cfg.AWSConfig.Endpoint) > 0 {
		c.Endpoint = aws.String(cfg.AWSConfig.Endpoint)
	}
	// Only set credentials if key and secret are given, otherwise fall back to IAM role
	if len(cfg.AWSConfig.Key) > 0 && len(cfg.AWSConfig.Secret) > 0 {
		c.Credentials = credentials.NewStaticCredentials(cfg.AWSConfig.Key, cfg.AWSConfig.Secret, "")
	}
	return c
}

// FOR TESTING

func localRoute53Config() *config.Config {
	return &config.Config{AWSConfig: config.AWSConfig{
		Endpoint:        "http://localhost:4580",
		Region:          "us-east-1",
		Key:             "123",
		Secret:          "456",
		DynamoTableName: TESTING_TABLE_NAME,
	}}
}

func localDynamoDBConfig(t *testing.T) *config.Config {
	cfg := &config.Config{AWSConfig: config.AWSConfig{
		Endpoint:        "http://localhost:4569",
		Region:          "us-east-1",
		Key:             "123",
		Secret:          "456",
		DynamoTableName: TESTING_TABLE_NAME,
	}}

	// for testing, let's validate that our table is created and create as needed:
	exists, _ := TableExists(context.Background(), cfg)
	if !exists {
		log.Println("Creating table: " + cfg.AWSConfig.DynamoTableName)
		_, err := CreateTable(context.Background(), cfg)
		ensure.Nil(t, err)
	}

	return cfg
}
