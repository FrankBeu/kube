// Package certmanager provides certmanager-crds and the default certmanager-installation
package certmanager

import (
	helm "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/kubeConfig"
	"thesym.site/kube/lib/namespace"
)

var (
	namespaceName        = "cert-manager"
	namespaceCertmanager = &namespace.Namespace{
		Name: namespaceName,
		Tier: namespace.NamespaceTierStructural,
	}
)

func CreateCertmanager(ctx *pulumi.Context) error {
	// https://cert-manager.io/docs/installation/helm/
	// https://artifacthub.io/packages/helm/cert-manager/cert-manager

	_, err := namespace.CreateNamespace(ctx, namespaceCertmanager)
	if err != nil {
		return err
	}

	_, err = certificate.CreateClusterIssuer(ctx, kubeConfig.DomainClusterIssuer(ctx))
	if err != nil {
		return err
	}

	_, err = helm.NewRelease(ctx, "cert-manager", &helm.ReleaseArgs{
		Version:   pulumi.String("1.11.0"),
		Chart:     pulumi.String("cert-manager"),
		Namespace: pulumi.String(namespaceCertmanager.Name),
		Values: pulumi.Map{
			"installCRDs": pulumi.Bool(true),
			//// GatewayAPI
			"extraArgs": pulumi.Array{
				pulumi.String("--feature-gates=ExperimentalGatewayAPISupport=true"),
			},
		},
		RepositoryOpts: &helm.RepositoryOptsArgs{
			Repo: pulumi.String("https://charts.jetstack.io"),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
