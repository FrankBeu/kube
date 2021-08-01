// Package loki installs a loki instance via helm-chart
package loki

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/kubeConfig"
	"thesym.site/kube/lib/namespace"
)

//nolint:funlen
func CreateLoki(ctx *pulumi.Context) error {
	name := "loki"

	namespaceLoki := &namespace.Namespace{
		Name: "loki",
		Tier: namespace.NamespaceTierMonitoring,
	}

	domainNameSuffix := kubeConfig.DomainNameSuffix(ctx)

	_, err := namespace.CreateNamespace(ctx, namespaceLoki)
	if err != nil {
		return err
	}
	_, err = helm.NewChart(ctx, name, helm.ChartArgs{
		// helm repo add grafana https://grafana.github.io/helm-charts
		Repo:      pulumi.String("grafana"),
		Chart:     pulumi.String("loki-stack"),
		Version:   pulumi.String("2.4.1"),
		Namespace: pulumi.String(namespaceLoki.Name),

		Values: pulumi.Map{
			"grafana": pulumi.Map{
				"enabled":       pulumi.Bool(true),
				"adminPassword": pulumi.String(kubeConfig.GrafanaAdminPassword(ctx)),
				"ingress": pulumi.Map{
					"enabled": pulumi.Bool(true),
					"annotations": pulumi.Map{
						"cert-manager.io/cluster-issuer": pulumi.String(certificate.ClusterIssuerTypeCaLocal.String()),
					},
					"ingressClassName": pulumi.String("nginx"),
					"hosts": pulumi.Array{
						pulumi.String("lokigraf" + domainNameSuffix),
					},
					"tls": pulumi.Array{
						pulumi.Map{
							"hosts": pulumi.Array{
								pulumi.String("lokigraf" + domainNameSuffix),
							},
							"secretName": pulumi.String("loki" + "-tls"),
						},
					},
				},
			},
			"prometheus": pulumi.Map{
				"enabled": pulumi.Bool(false),
				"alertmanager": pulumi.Map{
					"persistentVolume": pulumi.Map{
						"enabled": pulumi.Bool(false),
					},
				},
				"server": pulumi.Map{
					"persistentVolume": pulumi.Map{
						"enabled": pulumi.Bool(false),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
