package test

import (
	"testing"

	dnsassertions "github.com/armakuni/go-dns-assertions"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func toTerraformOptions(path string, vars *map[string]interface{}) terraform.Options {
	return terraform.Options{
		TerraformDir: path,
		Vars:         *vars,
	}
}

func TestTerraformAwsRoute53Zone(t *testing.T) {
	/* ARRANGE */
	options := toTerraformOptions("../../examples/complete", &map[string]interface{}{
		"zone_name": "terraform-test.armakuni.com.",
		"records": []map[string]interface{}{
			{"name": "one", "type": "A", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
			{"name": "two", "type": "TXT", "records": []string{"example-text-record"}, "ttl": 60},
		},
	})
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	///* ACTION */
	terraform.InitAndPlan(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	/* ASSERTIONS */
	zoneID := terraform.Output(t, terraformOptions, "zone_id")

	nameServers := GetRoute53HostedZoneNameServers(t, zoneID)
	if len(nameServers) < 1 {
		t.Errorf("No nameservers return for hosted zone")
		return
	}

	dnsServer := nameServers[0]

	dnsClient := dnsassertions.NewTestClient(t)
	lookupOne := dnsClient.FetchDNSRecords("one.terraform-test.armakuni.com", dnsServer)
	lookupOne.AssertHasARecord("10.0.0.0")
	lookupOne.AssertHasARecord("192.0.0.0")
	lookupOne.AssertHasTXTRecord("example-text-record")

	lookupTwo := dnsClient.FetchDNSRecords("two.terraform-test.armakuni.com", dnsServer)
	lookupTwo.AssertHasCNAMERecord("dummy.armakuni.co.uk.")
}
