package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"testing"
)

func GetRoute53HostedZoneNameServers(t *testing.T, zoneID string) []string {
	route53Session := session.Must(session.NewSession())

	var result []string

	svc := route53.New(route53Session)
	actualHostedZone, err := svc.GetHostedZone(&route53.GetHostedZoneInput{Id: aws.String(zoneID)})
	if err != nil {
		t.Errorf("Unable to GetHostedZone, %v", err)
		return result
	}

	if len((*actualHostedZone).DelegationSet.NameServers) < 1 {
		t.Errorf("No nameservers return for hosted zone")
		return result
	}

	for _, str := range (*actualHostedZone).DelegationSet.NameServers {
		result = append(result, *str)
	}

	return result
}
