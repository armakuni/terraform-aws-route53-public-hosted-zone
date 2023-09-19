package test

import (
	"os"
	"strings"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/require"
)

// We needed this new function to return the error
func InitAndPlanAndShowWithStructNoLogTempPlanFileE(t testing.TestingT, options *terraform.Options) (*terraform.PlanStruct, error) {
	oldLogger := options.Logger
	options.Logger = logger.Discard
	defer func() { options.Logger = oldLogger }()

	tmpFile, err := os.CreateTemp("", "terratest-plan-file-")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())
	defer require.NoError(t, os.Remove(tmpFile.Name()))

	options.PlanFilePath = tmpFile.Name()
	return terraform.InitAndPlanAndShowWithStructE(t, options)
}

func TFResourceType(resource string) string {
	resourceSegments := strings.Split(resource, ".")
	return resourceSegments[len(resourceSegments)-2]
}

func GetResourceChangeByAddress(address string, plan *terraform.PlanStruct) *tfjson.ResourceChange {
	for _, value := range plan.ResourceChangesMap {
		if value.Address == address {
			return value
		}
	}
	return nil
}

func GetResourceChangeAfter(resourceChange *tfjson.ResourceChange) map[string]interface{} {
	return resourceChange.Change.After.(map[string]interface{})
}

func interfaceSliceToStringSlice(input []interface{}) []string {
	result := make([]string, len(input))
	for i, val := range input {
		if strVal, ok := val.(string); ok {
			result[i] = strVal
		}
	}
	return result
}
