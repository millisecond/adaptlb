package awsutil

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/millisecond/adaptlb/config"
)

const DEFAULT_TTL = 60

func Route53Client(ctx context.Context, cfg *config.Config) *route53.Route53 {
	return route53.New(CreateSession(cfg))
}

func ListZones(ctx context.Context, cfg *config.Config) (*route53.ListHostedZonesOutput, error) {
	return Route53Client(ctx, cfg).ListHostedZones(&route53.ListHostedZonesInput{})
}

func CreateZone(ctx context.Context, cfg *config.Config, zoneName string) (*route53.CreateHostedZoneOutput, error) {
	return Route53Client(ctx, cfg).CreateHostedZone(&route53.CreateHostedZoneInput{
		CallerReference:  aws.String(zoneName + "create"),
		HostedZoneConfig: &route53.HostedZoneConfig{Comment: aws.String("LineSpeedLB Managed Zone")},
		Name:             aws.String(zoneName),
	})
}

func CurrentRecords(ctx context.Context, cfg *config.Config, zoneID string) ([]string, error) {
	output, err := Route53Client(ctx, cfg).ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
	})
	if err != nil {
		return nil, err
	}
	ips := []string{}
	for _, recordSet := range output.ResourceRecordSets {
		if *recordSet.Type != "A" {
			return nil, errors.New("Invalid record set type: " + *recordSet.Type)
		}
		for _, record := range recordSet.ResourceRecords {
			ips = append(ips, *record.Value)
		}
	}
	return ips, nil
}

func UpdateZone(ctx context.Context, cfg *config.Config, zoneID string, recordSetName string, targetIPs []string) (*route53.ChangeResourceRecordSetsOutput, error) {
	records := []*route53.ResourceRecord{}
	for _, ip := range targetIPs {
		records = append(records, &route53.ResourceRecord{
			Value: aws.String(ip),
		})
	}
	return Route53Client(ctx, cfg).ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String("LineSpeedLB Update"),
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Type:            aws.String("A"),
						TTL:             aws.Int64(DEFAULT_TTL),
						Name:            aws.String(recordSetName),
						ResourceRecords: records,
					},
				},
			},
		},
	})
}

func DeleteZone(ctx context.Context, cfg *config.Config, zoneID string) (*route53.DeleteHostedZoneOutput, error) {
	return Route53Client(ctx, cfg).DeleteHostedZone(&route53.DeleteHostedZoneInput{
		Id: aws.String(zoneID),
	})
}
