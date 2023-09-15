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
		"zone_name": "example.com.",
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
	route53ExpectedData := map[string]interface{}{
		"zone_name": "example.armakuni.com.",
		"records": []map[string]interface{}{
			{"name": "one", "type": "A", "records": []string{"10.0.0.0", "192.0.0.0"}, "ttl": 60},
			{"name": "two", "type": "CNAME", "records": []string{"dummy.armakuni.co.uk"}, "ttl": 60},
		},
	}

	options := toTerraformOptions("../../examples/complete", route53ExpectedData)
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &options)

	/* ACTION */
	plan, err := InitAndPlanAndShowWithStructNoLogTempPlanFileE(t, terraformOptions)
	if err != nil {
		t.Fatal(err)
	}

	/* ASSERTION */
	for resource, resourcePlan := range plan.ResourceChangesMap {
		switch TFResourceType(resource) {
		// These tests are only appropriate for a route53 record resource, skip other resource types
		case "aws_route53_record":
			// Record.Change.After is interface{} from raw terraform deserialised JSON
			recordMap, _ := resourcePlan.Change.After.(map[string]interface{})
			// Records is a list of strings
			recordsInterface := recordMap["records"].([]interface{})

			switch recordMap["type"] {
			case "A":
				// Check correct amount of records expected
				if len(recordsInterface) != 2 {
					t.Fatalf("Invalid 'A' records: %+v, should be %+v\n", recordsInterface, 2)
				}

				// Check each records matches expected
				for _, record := range recordsInterface {
					if recordStr, ok := record.(string); ok {
						fmt.Printf("Record: %+v\n", recordStr)
						t.Fail()
					}
				}
			case "CNAME":
				// Check each records matches expected
				for _, cname := range recordsInterface {
					if cnameStr, ok := cname.(string); ok {
						fmt.Printf("CNAME: %+v\n", cnameStr)
						t.Fail()
					}
				}
			default:
				t.Fatalf("Tests not implemented for record type: %s", recordMap["type"])
			}
		case "aws_route53_zone":
			zoneMap, _ := resourcePlan.Change.After.(map[string]interface{})
			assert.EqualValues(t, route53ExpectedData["zone_name"], zoneMap["name"])
		default:
			t.Fatalf("Tests not implemented for resource type: %s", TFResourceType(resource))
		}
	}

	// Number of Resources to be created
	assert.EqualValues(t, 3, len(plan.ResourcePlannedValuesMap))
}
