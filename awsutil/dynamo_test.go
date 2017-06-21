package awsutil

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/facebookgo/ensure"
	"testing"
)

func TestDyanmoVersioning(t *testing.T) {
	cfg := localDynamoDBConfig(t)

	DeleteItem(context.Background(), cfg, DNSUpdateType, DNSUpdateRangeKey)

	_, version, err := GetVersionedItem(context.Background(), cfg, DNSUpdateType, DNSUpdateRangeKey)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, version, -1)

	values := map[string]*dynamodb.AttributeValue{
		"records": {SS: aws.StringSlice([]string{"1.1.1.1", "1.1.1.2"})},
	}

	updated, err := UpdateVersionedItem(context.Background(), cfg, DNSUpdateType, DNSUpdateRangeKey, version, values)
	ensure.Nil(t, err)
	ensure.True(t, updated)

	// update again to same version should fail
	updated, err = UpdateVersionedItem(context.Background(), cfg, DNSUpdateType, DNSUpdateRangeKey, version, values)
	ensure.Nil(t, err)
	ensure.False(t, updated)

	ips := []string{"1.1.1.1", "1.1.1.2", "1.1.1.3"}
	values = map[string]*dynamodb.AttributeValue{
		"records": {SS: aws.StringSlice(ips)},
	}

	updated, err = UpdateVersionedItem(context.Background(), cfg, DNSUpdateType, DNSUpdateRangeKey, version+1, values)
	ensure.Nil(t, err)
	ensure.True(t, updated)

	get, err := GetItem(context.Background(), cfg, DNSUpdateType, DNSUpdateRangeKey)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, aws.StringValueSlice(get.Item["records"].SS), ips)
}
