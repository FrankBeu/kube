// Package gitea installs a gitea instance via helm-chart
package gitea

import (
	"strings"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"thesym.site/kube/lib/kubeconfig"
	"thesym.site/kube/lib/namespace"
	"thesym.site/kube/lib/persistence"
	"thesym.site/kube/lib/types"
)

type giteaSecret struct {
	PostgresqlPassword string
	PostgresqlUsername string
}

var (
	name = "gitea"
	// subDomainName = "git"
	subDomainName = name

	namespaceGitea = &namespace.Namespace{
		Name: name,
		Tier: namespace.NamespaceTierCommunication,
	}

	persistentDataConfig = &persistence.Config{
		Name:                  "gitea-data",
		StorageClass:          persistence.StorageClassLocalPath,
		Capacity:              "10Gi",
		AccessMode:            persistence.AccessModeRWO,
		PersistencePathSuffix: "gitea/data",
		NamespaceName:         namespaceGitea.Name,
	}
	//// claim will be created via helm
	persistentDBConfig = &persistence.Config{
		Name:                  "data-gitea-postgresql-0",
		StorageClass:          persistence.StorageClassLocalPath,
		Capacity:              "10Gi",
		AccessMode:            persistence.AccessModeRWO,
		PersistencePathSuffix: "gitea/db",
		NamespaceName:         namespaceGitea.Name,
	}
)

// CreateGitea creates a gitea instance via helm
//
//nolint:funlen
func CreateGitea(ctx *pulumi.Context) error {
	conf := config.New(ctx, "")
	// domainName := subDomainName + "." + conf.Require("domain")
	domainName := subDomainName + kubeconfig.DomainNameSuffix(ctx)

	_, err := namespace.CreateNamespace(ctx, namespaceGitea)
	if err != nil {
		return err
	}

	err = createPersistence(ctx)
	if err != nil {
		return err
	}

	var gs giteaSecret
	conf.RequireSecretObject("giteaSecret", &gs)
	postgresqlPassword := pulumi.String(gs.PostgresqlPassword)
	postgresqlUsername := pulumi.String(gs.PostgresqlUsername)

	_ = postgresqlPassword
	_ = postgresqlUsername

	_, err = helm.NewChart(ctx, "gitea", helm.ChartArgs{
		// `helm repo add gitea-charts https://dl.gitea.io/charts/`
		Repo:      pulumi.String("gitea-charts"),
		Chart:     pulumi.String("gitea"),
		Version:   pulumi.String("3.1.4"),
		Namespace: pulumi.String(namespaceGitea.Name),
		//// APIVersion is NOT working: https://github.com/pulumi/pulumi-kubernetes/issues/1034
		//// use Transformations instead
		APIVersions: pulumi.StringArray{pulumi.String("networking.k8s.io/v1")},
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
				"annotations": pulumi.Map{
					"cert-manager.io/cluster-issuer": pulumi.String(types.ClusterIssuerTypeCALocal.String()),
				},
				"hosts": pulumi.Array{
					pulumi.String(domainName),
				},
				"tls": pulumi.Array{
					pulumi.Map{
						"secretName": pulumi.String(name + "-tls"),
						"hosts": pulumi.Array{
							pulumi.String(domainName),
						},
					},
				},
			},
			"resources": pulumi.Map{
				"requests": pulumi.Map{
					"cpu": pulumi.String("10m"), // resources.requests.cpu=10m
				},
			},
			"persistence": pulumi.Map{
				"existingClaim": pulumi.String(persistentDataConfig.Name),
			},
			"postgresql": pulumi.Map{
				"global": pulumi.Map{
					"postgresql": pulumi.Map{
						"postgresqlUsername": postgresqlUsername,
						"postgresqlPassword": postgresqlPassword,
					},
				},
			},
		},
		Transformations: []yaml.Transformation{
			func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
				//// only act if the ingress is passed to this callback
				if state["kind"] == "Ingress" {
					state["apiVersion"] = "networking.k8s.io/v1"

					//// backendInformation{Retrieval,Setting}
					//nolint:lll
					paths := state["spec"].(map[string]interface{})["rules"].([]interface{})[0].(map[string]interface{})["http"].(map[string]interface{})["paths"]
					path := paths.([]interface{})[0].(map[string]interface{})["path"]
					backend := paths.([]interface{})[0].(map[string]interface{})["backend"]
					serviceName := backend.(map[string]interface{})["serviceName"]
					servicePort := backend.(map[string]interface{})["servicePort"]

					//// specRebuild
					spec := map[string]interface{}{
						"tls": []interface{}{
							map[string]interface{}{
								"hosts": []interface{}{
									domainName,
								},
								"secretName": name + "-tls",
							},
						},
						"ingressClassName": "nginx",
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
	})
	if err != nil {
		return err
	}

	return nil
}

func createPersistence(ctx *pulumi.Context) error {
	_, err := persistence.CreatePersistentVolume(ctx, persistentDataConfig)
	if err != nil {
		return err
	}

	_, err = persistence.CreatePersistentVolumeClaim(ctx, persistentDataConfig)
	if err != nil {
		return err
	}

	_, err = persistence.CreatePersistentVolume(ctx, persistentDBConfig)
	if err != nil {
		return err
	}

	return nil
}
