package test

import (
	"fmt"
	"testing"

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
	options := toTerraformOptions("../../examples/complete", map[string]interface{}{
		"zone_name": "example.com",
		"records": []map[string]interface{}{
			{"name": "one", "type": "ALPHA", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
		},
	})
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	_, err := InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)

	/* ASSERTION */
	assert.ErrorContains(t, err, `Only valid types permitted (A, CNAME, MX, NS, TXT, SOA, SPF)`)
}

func TestRoute53ZoneHasValidRecordEntries(t *testing.T) {
	/* ARRANGE */
	inputVariables := map[string]interface{}{
		"zone_name": "example.armakuni.com",
		"records": []map[string]interface{}{
			{"name": "one", "type": "A", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
			{"name": "three", "type": "CNAME", "records": []string{"mummy.armakuni.co.uk"}, "ttl": 60},
		},
	}

	options := toTerraformOptions("../../examples/complete", inputVariables)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	plan, err := InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)
	assert.Empty(t, err)

	/* ASSERTIONS */
	// Asserting Zone Name
	actualZone := GetResourceChangeByAddress("module.test_route53_zone.aws_route53_zone.this", plan)
	assert.EqualValues(t, inputVariables["zone_name"], GetResourceChangeAfter(actualZone)["name"])
	// Asserting Records
	for _, expectedValue := range inputVariables["records"].([]map[string]interface{}) {
		address := fmt.Sprintf("module.test_route53_zone.aws_route53_record.record[\"name=%s,type=%s\"]", expectedValue["name"], expectedValue["type"])
		actualRecordResourceChange := GetResourceChangeByAddress(address, plan)
		assert.NotEmpty(t, actualRecordResourceChange, fmt.Sprintf("ResourceChange for: %s does not exist", address))

		actualRecordResourceChangeAfter := GetResourceChangeAfter(actualRecordResourceChange)
		expectedRecordName := fmt.Sprintf("%s.%s", expectedValue["name"], inputVariables["zone_name"])
		assert.EqualValues(t, expectedRecordName, actualRecordResourceChangeAfter["name"])
		assert.EqualValues(t, expectedValue["type"], actualRecordResourceChangeAfter["type"])
		assert.EqualValues(t, expectedValue["ttl"], actualRecordResourceChangeAfter["ttl"])
		actualRecords := interfaceSliceToStringSlice(actualRecordResourceChangeAfter["records"].([]interface{}))
		assert.EqualValues(t, expectedValue["records"], actualRecords)
	}
}
