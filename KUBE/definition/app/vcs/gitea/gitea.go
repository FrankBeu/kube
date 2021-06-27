// Package giteas installs a gitea instance via helm-chart
package gitea

import (
	"strings"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"thesym.site/kube/lib"
)

var (
	name = "gitea"

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
	conf := config.New(ctx, "")
	domainName := name + "." + conf.Require("domain")

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
		Transformations: []yaml.Transformation{
			// Make every service private to the cluster, i.e., turn all services into ClusterIP
			// instead of LoadBalancer.
			func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
				// name := state["metadata"].(map[string]interface{})["name"]
				// if state["kind"] == "Pod" && name == "test" {
				if state["kind"] == "Ingress" {
					state["apiVersion"] = "networking.k8s.io/v1"
				}
				// spec := state["spec"].(map[string]interface{})
				// // if spec["
				// spec["type"] = "ClusterIP"
			},
		},

		//// Do not act on the yamlLayer - act on the helmValuesLayer
		//// DEBUG: NOT: `helm template -s templates/gitea/ingress.yaml ./CHART/gitea --set ingress.enabled=true --set "ingress.hosts\.0.host"=git.thesym.site`
		//// DEBUG: `cat CHART/gitea/templates/gitea/ingress.yaml`
		Values: pulumi.Map{
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
			"ingress": pulumi.Map{
				"enabled": pulumi.Bool(true),
				// "annotations": pulumi.Map{
				// 	"": pulumi.String(""),
				// },
				"hosts": pulumi.Array{ // ingress.hosts[0].host
					pulumi.String(domainName),
				},

				// --set resources.requests.cpu=10m \
				// --set persistence.existingClaim={{.pvName0}} \
				// --set gitea.config.APP_NAME=Gitea \
				// --set gitea.config.server.DOMAIN={{.domainName}} \
				// --set gitea.config.server.ROOT_URL=https://{{.domainName}} \
				// --set gitea.config.server.SSH_DOMAIN={{.domainName}} \

			},
			"resources": pulumi.Map{
				"requests": pulumi.Map{
					"cpu": pulumi.String("10m"),
				},
			},
			// --set persistence.existingClaim={{.pvName0}} \
			// "persistence": pulumi.Map{
			// "existingClaim": pulumi.String(pvNameData),
			// },
		},
		Namespace: pulumi.String(namespaceGitea.Name),
	})
	if err != nil {
		return err
	}

	return nil
}
