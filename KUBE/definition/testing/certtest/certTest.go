// Package certtest use used for testing certManager
package certtest

import (
	"strconv"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib"
)

var (
	exampleCert = lib.Cert{
		ClusterIssuerType: lib.ClusterIssuerTypeCaLocal,
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
	err := lib.CreateCert(ctx, &exampleCert)
	if err != nil {
		return err
	}

	return nil
}
