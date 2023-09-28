package test

import (
	"fmt"
	"testing"

	"github.com/armakuni/go-terratest-helper"
	"github.com/armakuni/go-terratest-helper/tfplan"
	"github.com/armakuni/go-terratest-helper/utils"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func toTerraformOptions(path string, vars map[string]interface{}) terraform.Options {
	return terraform.Options{
		TerraformDir: path,
		Vars:         vars,
	}
}

func TestRoute53ZoneWhenInvalidRecordTypeIsPassed(t *testing.T) {
	/* ARRANGE */
	t.Parallel()
	options := toTerraformOptions("../../examples/complete", map[string]interface{}{
		"zone_name": "example.com",
		"records": []map[string]interface{}{
			{"name": "one", "type": "ALPHA", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
		},
	})
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	_, err := tfplan.InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)

	/* ASSERTION */
	assert.ErrorContains(t, err, `Only valid types permitted (A, CNAME, MX, NS, TXT, SOA, SPF)`)
}

func TestRoute53ZoneHasValidRecordEntries(t *testing.T) {
	/* ARRANGE */
	t.Parallel()
	inputVariables := map[string]interface{}{
		"zone_name": "example.armakuni.com",
		"records": []map[string]interface{}{
			{"name": "one", "type": "A", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
			{"name": "three", "type": "TXT", "records": []string{"example-text-record"}, "ttl": 60},
			{"name": "four", "type": "NS", "records": []string{"ns-2121.awsdns-21.com."}, "ttl": 60},
			{"name": "five", "type": "SOA", "records": []string{"ns-2121.awsdns-21.com."}, "ttl": 60},
			{"name": "six", "type": "SPF", "records": []string{"v=spf1 ip4:10.10.10.10/16-all"}, "ttl": 60},
			{"name": "seven", "type": "MX", "records": []string{"10 mailserver.dummy.armakuni.co.uk"}, "ttl": 60},
		},
	}

	options := toTerraformOptions("../../examples/complete", inputVariables)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	plan, err := tfplan.InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)
	assert.Empty(t, err)

	/* ASSERTIONS */
	// Asserting Zone Name
	actualZone, _ := tfplanstruct.GetResourceChangeAfterByAddressE("module.test_route53_zone.aws_route53_zone.this", plan)
	assert.EqualValues(t, inputVariables["zone_name"], actualZone["name"])
	// Asserting Records, Iterating over the expected Route53 Records and asserting the against the plan
	for _, expectedValue := range inputVariables["records"].([]map[string]interface{}) {
		recordTFPlanAddress := fmt.Sprintf("module.test_route53_zone.aws_route53_record.record[\"name=%s,type=%s\"]", expectedValue["name"], expectedValue["type"])
		actualRecordResourceChangeAfter, _ := tfplanstruct.GetResourceChangeAfterByAddressE(recordTFPlanAddress, plan)
		assert.NotEmpty(t, actualRecordResourceChangeAfter, fmt.Sprintf("ResourceChange for: %s does not exist", recordTFPlanAddress))

		expectedRecordName := fmt.Sprintf("%s.%s", expectedValue["name"], inputVariables["zone_name"])
		assert.EqualValues(t, expectedRecordName, actualRecordResourceChangeAfter["name"])
		assert.EqualValues(t, expectedValue["type"], actualRecordResourceChangeAfter["type"])
		assert.EqualValues(t, expectedValue["ttl"], actualRecordResourceChangeAfter["ttl"])
		actualRecordsArray, _ := utils.InterfaceSliceToStringSliceE(actualRecordResourceChangeAfter["records"].([]interface{}))
		assert.EqualValues(t, expectedValue["records"], actualRecordsArray)
	}
}
