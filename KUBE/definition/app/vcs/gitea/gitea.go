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
	name          = "gitea"
	subDomainName = "git"

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
	domainName := subDomainName + "." + conf.Require("domain")

	err := lib.CreateNamespaces(ctx, namespaceGitea)
	if err != nil {
		return err
	}
	// `helm repo add gitea-charts https://dl.gitea.io/charts/`
	_, err = helm.NewChart(ctx, "gitea", helm.ChartArgs{
		Repo:      pulumi.String("gitea-charts"),
		Chart:     pulumi.String("gitea"),
		Version:   pulumi.String("3.1.4"),
		Namespace: pulumi.String(namespaceGitea.Name),
		//// APIVersion is NOT working: https://github.com/pulumi/pulumi-kubernetes/issues/1034
		//// use Trasnformations instead
		APIVersions: pulumi.StringArray{pulumi.String("networking.k8s.io/v1")},
		// TODO: test if `helm repo add ...` is necessary
		// FetchArgs: &helm.FetchArgs{
		// Repo: pulumi.String("https://dl.gitea.io/charts/"),
		// },
		Transformations: []yaml.Transformation{
			func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
				//// only act if the ingress is passed to this callback
				if state["kind"] == "Ingress" {
					state["apiVersion"] = "networking.k8s.io/v1"

					//// information{Retrieval,Setting}
					//nolint:lll
					paths := state["spec"].(map[string]interface{})["rules"].([]interface{})[0].(map[string]interface{})["http"].(map[string]interface{})["paths"]
					path := paths.([]interface{})[0].(map[string]interface{})["path"]
					backend := paths.([]interface{})[0].(map[string]interface{})["backend"]
					serviceName := backend.(map[string]interface{})["serviceName"]
					servicePort := backend.(map[string]interface{})["servicePort"]

					//// specRebuild
					spec := map[string]interface{}{
						"rules": []interface{}{
							map[string]interface{}{
								"host": domainName,
								"http": map[string]interface{}{
									"paths": []interface{}{
										map[string]interface{}{
											"path":     path,
											"pathType": "Prefix",
											"backend": map[string]interface{}{
												"service": map[string]interface{}{
													"name": serviceName,
													"port": map[string]interface{}{
														"number": servicePort,
													},
												},
											},
										},
									},
								},
							},
						},
					}
					state["spec"] = spec
				}
			},
		},
		Values: pulumi.Map{
			"gitea": pulumi.Map{
				"config": pulumi.Map{
					"APP_NAME": pulumi.String(strings.Title(name)),
					"server": pulumi.Map{
						"DOMAIN":     pulumi.String(domainName),
						"ROOT_URL":   pulumi.String("https://" + domainName),
						"SSH_DOMAIN": pulumi.String(domainName), // gitea.config.server.SSH_DOMAIN={{.domainName}}
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
			},
			"resources": pulumi.Map{
				"requests": pulumi.Map{
					"cpu": pulumi.String("10m"), // resources.requests.cpu=10m
				},
			},
			// --set persistence.existingClaim={{.pvName0}} \
			// "persistence": pulumi.Map{
			// "existingClaim": pulumi.String(pvNameData),
			// },
		},
	})
	if err != nil {
		return err
	}

	return nil
}
