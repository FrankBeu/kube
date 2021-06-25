// Package giteas installs a gitea instance via helm-chart
package gitea

import (
	"strings"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/lib"
)

var (
	name       = "gitea"
	domainName = name + ".thesym.site"
	// pvNameData = name + "-data"

	// pvName0:                  "{{.name}}-data"
	// pvName1:                  "data-{{.name}}-postgresql-0"
	// pvRelativePath0:          "{{.name}}/data"
	// pvRelativePath1:          "{{.name}}/db"
	// pvStorageSize0:           10Gi
	// pvStorageSize1:           10Gi
	// pvStorageClass0:          local-path
	// pvStorageClass1:          local-path

	namespaceGitea = &lib.Namespace{
		Name: name,
		Tier: lib.NamespaceTierCommunication,
		// GlooDiscovery: true,
	}
)

func CreateGitea(ctx *pulumi.Context) error {

	err := lib.CreateNamespaces(ctx, namespaceGitea)
	if err != nil {
		return err
	}
	// `helm repo add gitea-charts https://dl.gitea.io/charts/`
	_, err = helm.NewChart(ctx, "gitea", helm.ChartArgs{
		Repo:    pulumi.String("gitea-charts"),
		Chart:   pulumi.String("gitea"),
		Version: pulumi.String("3.1.4"),
		// TODO: test if `helm repo add ...` is necessary
		// FetchArgs: &helm.FetchArgs{
		// Repo: pulumi.String("https://dl.gitea.io/charts/"),
		// },

		//// Do not act on the yamlLayer - act on the helmValuesLayer
		//// DEBUG: NOT: `helm template -s templates/gitea/ingress.yaml ./CHART/gitea --set ingress.enabled=true --set "ingress.hosts\.0.host"=git.thesym.site`
		//// DEBUG: `cat CHART/gitea/templates/gitea/ingress.yaml`
		Values: pulumi.Map{
			"ingress": pulumi.Map{
				"enabled": pulumi.Bool(true),
				// "annotations": pulumi.Map{
				// 	"": pulumi.String(""),
				// },
				"gitea": pulumi.Map{
					"config": pulumi.Map{
						"APP_NAME": pulumi.String(strings.Title(name)),
						"server": pulumi.Map{

							"DOMAIN":     pulumi.String(domainName),
							"ROOT_URL":   pulumi.String("https://" + domainName),
							"SSH_DOMAIN": pulumi.String(domainName),
						},
					},
				},
				"hosts": pulumi.Array{ // ingress.hosts[0].host
					pulumi.String(domainName),
				},

				// --set persistence.existingClaim={{.pvName0}} \
				// "persistence": pulumi.Map{
				// "existingClaim": pulumi.String(pvNameData),
				// },
				"resources": pulumi.Map{
					"requests": pulumi.Map{
						"cpu": pulumi.String("10m"),
					},
				},

				// --set resources.requests.cpu=10m \
				// --set persistence.existingClaim={{.pvName0}} \
				// --set gitea.config.APP_NAME=Gitea \
				// --set gitea.config.server.DOMAIN={{.domainName}} \
				// --set gitea.config.server.ROOT_URL=https://{{.domainName}} \
				// --set gitea.config.server.SSH_DOMAIN={{.domainName}} \

			},
		},
		Namespace: pulumi.String(namespaceGitea.Name),
	})
	if err != nil {
		return err
	}

	return nil
}
