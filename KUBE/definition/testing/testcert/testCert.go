// Package testcert use used for testing certManager
package testcert

import (
	"strconv"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/types"
)

var (
	exampleCert = types.Cert{
		ClusterIssuerType: types.ClusterIssuerTypeCALocal,
		Namespace:         "test",
		Name:              "testing",
		Duration:          strconv.Itoa(99*24) + "h",
		AdditionalSubdomainNames: []string{
			"tests",
			"test",
		},
	}
)

// CreateTestCert creates an exampleCertificate
func CreateTestCert(ctx *pulumi.Context) error {
	err := certificate.CreateCert(ctx, &exampleCert)
	if err != nil {
		return err
	}

	return nil
}
