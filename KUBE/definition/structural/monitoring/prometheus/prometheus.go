// Package prometheus installs a prometheus instance via helm-chart
package prometheus

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/namespace"
)

type PromStackConfig struct {
	Grafana struct {
		AdminPassword string
	}
}

//nolint:funlen
func CreatePrometheus(ctx *pulumi.Context) error {
	name := "prometheus"

	namespacePrometheus := &namespace.Namespace{
		Name: "monitoring",
		Tier: namespace.NamespaceTierCommunication,
		// GlooDiscovery: true,
	}

	conf := config.New(ctx, "")

	domainNameSuffix := "." + conf.Require("domain")

	var psc PromStackConfig
	conf.RequireSecretObject("promStack", &psc)
	grafanaAdminPw := psc.Grafana.AdminPassword

	_, err := namespace.CreateNamespace(ctx, namespacePrometheus)
	if err != nil {
		return err
	}
	_, err = helm.NewChart(ctx, name, helm.ChartArgs{
		// helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
		Repo:      pulumi.String("prometheus-community"),
		Chart:     pulumi.String("kube-prometheus-stack"),
		Version:   pulumi.String("16.12.1"),
		Namespace: pulumi.String(namespacePrometheus.Name),
		Values: pulumi.Map{
			"defaultRules": pulumi.Map{
				"etcd":          pulumi.Bool(false),
				"kubeScheduler": pulumi.Bool(false),
			},
			"alertmanager": pulumi.Map{
				"ingress": pulumi.Map{
					"enabled": pulumi.Bool(true),
					"annotations": pulumi.Map{
						"cert-manager.io/cluster-issuer": pulumi.String(certificate.ClusterIssuerTypeCaLocal.String()),
					},
					"ingressClassName": pulumi.String("nginx"),
					"hosts": pulumi.Array{
						pulumi.String("alertmanager" + domainNameSuffix),
					},
					"paths": pulumi.Array{
						pulumi.String("/"),
					},
					"pathType": pulumi.String("ImplementationSpecific"),
					"tls": pulumi.Array{
						pulumi.Map{
							"hosts": pulumi.Array{
								pulumi.String("alertmanager" + domainNameSuffix),
							},
							"secretName": pulumi.String("alertmanager" + "-tls"),
						},
					},
				},
			},
			"grafana": pulumi.Map{
				"adminPassword": pulumi.String(grafanaAdminPw),
				"ingress": pulumi.Map{
					"enabled": pulumi.Bool(true),
					"annotations": pulumi.Map{
						"cert-manager.io/cluster-issuer": pulumi.String(certificate.ClusterIssuerTypeCaLocal.String()),
					},
					"ingressClassName": pulumi.String("nginx"),
					"hosts": pulumi.Array{
						pulumi.String("promgraf" + domainNameSuffix),
					},
					"tls": pulumi.Array{
						pulumi.Map{
							"hosts": pulumi.Array{
								pulumi.String("promgraf" + domainNameSuffix),
							},
							"secretName": pulumi.String("prometheus-grafana" + "-tls"),
						},
					},
				},
			},
			"kubeControllerManager": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"kubeEtcd": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"kubeProxy": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"kubeScheduler": pulumi.Map{
				"enabled": pulumi.Bool(false),
			},
			"prometheus": pulumi.Map{
				"ingress": pulumi.Map{
					"enabled": pulumi.Bool(true),
					"annotations": pulumi.Map{
						"cert-manager.io/cluster-issuer": pulumi.String(certificate.ClusterIssuerTypeCaLocal.String()),
					},
					"ingressClassName": pulumi.String("nginx"),
					"hosts": pulumi.Array{
						pulumi.String(name + domainNameSuffix),
					},
					"paths": pulumi.Array{
						pulumi.String("/"),
					},
					"pathType": pulumi.String("ImplementationSpecific"),
					"tls": pulumi.Array{
						pulumi.Map{
							"hosts": pulumi.Array{
								pulumi.String(name + domainNameSuffix),
							},
							"secretName": pulumi.String(name + "-tls"),
						},
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
