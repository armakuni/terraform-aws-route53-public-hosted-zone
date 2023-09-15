package test

import (
	"os"
	"strings"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
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
