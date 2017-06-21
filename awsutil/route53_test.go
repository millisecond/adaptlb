package awsutil

import (
	"context"
	"log"
	"github.com/facebookgo/ensure"
	"github.com/millisecond/adaptlb/config"
	"github.com/millisecond/adaptlb/testutil"
	"testing"
)

func TestCreateListDeleteZone(t *testing.T) {
	cfg := localRoute53Config()

	name := testutil.RandomString(10) + ".com."
	createOutput, err := CreateZone(context.Background(), cfg, name)
	ensure.Nil(t, err)
	log.Println(*createOutput.HostedZone.Id)

	removeZone(t, cfg, name)
}

func TestModifyZone(t *testing.T) {
	cfg := localRoute53Config()

	zoneName := testutil.RandomString(10) + ".com."
	createOutput, err := CreateZone(context.Background(), cfg, zoneName)
	ensure.Nil(t, err)
	log.Println(*createOutput.HostedZone.Id)

	updateOutput, err := UpdateZone(context.Background(), cfg, *createOutput.HostedZone.Id, zoneName, []string{"1.2.3.4", "1.2.3.5"})
	ensure.Nil(t, err)
	log.Println(updateOutput)

	current, err := CurrentRecords(context.Background(), cfg, *createOutput.HostedZone.Id)
	ensure.DeepEqual(t, len(current), 2)

	removeZone(t, cfg, zoneName)
}

func removeZone(t *testing.T, cfg *config.Config, name string) {
	found := false
	listOutput, err := ListZones(context.Background(), cfg)
	ensure.Nil(t, err)
	for _, zone := range listOutput.HostedZones {
		log.Println(zone)
		if *zone.Name == name {
			found = true
			_, err := DeleteZone(context.Background(), cfg, *zone.Id)
			ensure.Nil(t, err)
		}
	}
	ensure.True(t, found)

	found = false
	listOutput, err = ListZones(context.Background(), cfg)
	ensure.Nil(t, err)
	for _, zone := range listOutput.HostedZones {
		log.Println(zone)
		if *zone.Name == name {
			found = true
		}
	}
	ensure.False(t, found)
}
