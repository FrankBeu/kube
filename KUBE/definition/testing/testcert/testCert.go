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
		Duration:          strconv.Itoa(types.CertificateDefaultDurationInDays*24) + "h", //nolint:gomnd
		AdditionalSubdomainNames: []string{
			"tests",
			"test",
		},
	}
)

// CreateTestCert creates an exampleCertificate
func CreateTestCert(ctx *pulumi.Context) error {
	_, err := certificate.CreateCert(ctx, &exampleCert)
	if err != nil {
		return err
	}

	return nil
}
