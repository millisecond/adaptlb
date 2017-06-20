package config

type Config struct {
	AWSConfig AWSConfig
}

type AWSConfig struct {
	Region          string
	Endpoint        string
	DynamoTableName string
	Key             string
	Secret          string
}
