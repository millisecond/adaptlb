package awsutil

import ()
import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/millisecond/linespeedlb/config"
	"strconv"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

const DEFAULT_TABLE_NAME = "LineSpeedLB"
const TESTING_TABLE_NAME = "LineSpeedLBTesting"

const HASH_KEY = "h"
const RANGE_KEY = "r"
const VERSION_KEY = "v"

const DEFAULT_WRITE_THROUGHPUT = 5
const DEFAULT_READ_THROUGHPUT = 5

type ObjectType string

const (
	DNSUpdateType ObjectType = "DNSUPDATE"
	DNSUpdateRangeKey             = "dns"
)

func DynamoDBClient(ctx context.Context, cfg *config.Config) *dynamodb.DynamoDB {
	return dynamodb.New(CreateSession(cfg))
}

func TableExists(ctx context.Context, cfg *config.Config) (bool, error) {
	_, err := DynamoDBClient(ctx, cfg).DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(cfg.AWSConfig.DynamoTableName),
	})
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func CreateTable(ctx context.Context, cfg *config.Config) (*dynamodb.CreateTableOutput, error) {
	return DynamoDBClient(ctx, cfg).CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(cfg.AWSConfig.DynamoTableName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(HASH_KEY),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(RANGE_KEY),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(HASH_KEY),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String(RANGE_KEY),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(DEFAULT_READ_THROUGHPUT),
			WriteCapacityUnits: aws.Int64(DEFAULT_WRITE_THROUGHPUT),
		},
	})
}

func DeleteItem(ctx context.Context, cfg *config.Config, objectType ObjectType, id string) error {
	_, err := DynamoDBClient(ctx, cfg).DeleteItem(&dynamodb.DeleteItemInput{
		TableName:      aws.String(cfg.AWSConfig.DynamoTableName),
		Key: map[string]*dynamodb.AttributeValue{
			HASH_KEY:  {S: aws.String(string(objectType))},
			RANGE_KEY: {S: aws.String(id)},
		},
	})
	return err
}

func GetItem(ctx context.Context, cfg *config.Config, objectType ObjectType, id string) (*dynamodb.GetItemOutput, error) {
	return DynamoDBClient(ctx, cfg).GetItem(&dynamodb.GetItemInput{
		TableName:      aws.String(cfg.AWSConfig.DynamoTableName),
		Key: map[string]*dynamodb.AttributeValue{
			HASH_KEY:  {S: aws.String(string(objectType))},
			RANGE_KEY: {S: aws.String(id)},
		},
	})
}

func GetVersionedItem(ctx context.Context, cfg *config.Config, objectType ObjectType, id string) (map[string]*dynamodb.AttributeValue, int, error) {
	output, err := DynamoDBClient(ctx, cfg).GetItem(&dynamodb.GetItemInput{
		TableName:      aws.String(cfg.AWSConfig.DynamoTableName),
		ConsistentRead: aws.Bool(true),
		Key: map[string]*dynamodb.AttributeValue{
			HASH_KEY:  {S: aws.String(string(objectType))},
			RANGE_KEY: {S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, -1, err
	}
	if output == nil || output.Item == nil {
		// Not found
		return nil, -1, nil
	}
	version, err := strconv.Atoi(*output.Item[VERSION_KEY].N)
	if err != nil {
		return nil, -1, err
	}
	return output.Item, version, err
}

func UpdateVersionedItem(ctx context.Context, cfg *config.Config, objectType ObjectType, id string, oldVersion int, set map[string]*dynamodb.AttributeValue) (bool, error) {
	oldVersionStr := strconv.Itoa(oldVersion)
	expressionValues := map[string]*dynamodb.AttributeValue{}
	setStr := "SET " +VERSION_KEY +" = :" + VERSION_KEY +", "
	expressionValues[":" + VERSION_KEY] = &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(oldVersion + 1))}
	for k, v := range set {
		setStr += k + " = :" + k
		expressionValues[":"+k] = v
	}
	conditionalExpression := VERSION_KEY + " = :old" + VERSION_KEY
	if oldVersion < 0 {
		// no previous id, make sure we're the first one
		conditionalExpression = "attribute_not_exists("+VERSION_KEY+")"
	} else {
		expressionValues[":old" + VERSION_KEY] = &dynamodb.AttributeValue{N: aws.String(oldVersionStr)}
	}
	output, err := DynamoDBClient(ctx, cfg).UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(cfg.AWSConfig.DynamoTableName),
		Key: map[string]*dynamodb.AttributeValue{
			HASH_KEY:  {S: aws.String(string(objectType))},
			RANGE_KEY: {S: aws.String(id)},
		},
		ConditionExpression:       aws.String(conditionalExpression),
		ExpressionAttributeValues: expressionValues,
		UpdateExpression:          aws.String(setStr),
		ReturnValues:              aws.String("UPDATED_NEW"),
	})
	if e, ok := err.(awserr.Error); ok {
		if e.Code() == "ConditionalCheckFailedException" {
			return false, nil
		}
	}
	if err != nil {
		return false, err
	}
	// If we got back a version key, it was correctly updated
	_, contains := output.Attributes[VERSION_KEY]
	return contains, err
}
