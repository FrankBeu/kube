// Package jaeger installs a jaeger instance
package jaeger

import (
	"strings"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	jaegerv1 "thesym.site/kube/crds/jaeger/jaegertracing/v1"
	"thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/crd"
	"thesym.site/kube/lib/ingress"

	// "thesym.site/kube/lib/certificate"
	// "thesym.site/kube/lib/ingress"
	"thesym.site/kube/lib/namespace"
)

func CreateJaegerOperator(ctx *pulumi.Context) error {
	name := "jaeger"

	namespaceJaeger := &namespace.Namespace{
		Name: "observability",
		Tier: namespace.NamespaceTierMonitoring,
		// GlooDiscovery: true,
	}

	_, err := namespace.CreateNamespace(ctx, namespaceJaeger)
	if err != nil {
		return err
	}

	jaegerCrds := []crd.Crd{
		{Name: "jaeger-definition", Location: "./crds/jaeger/crdDefinitions/jaegertracing.io_jaegers_crd.yaml"},
	}
	err = crd.RegisterCRDs(ctx, jaegerCrds)
	if err != nil {
		return err
	}

	err = execGeneratedCode(ctx, namespaceJaeger)
	if err != nil {
		return err
	}

	err = createJaegerInstance(ctx, name, namespaceJaeger)

	if err != nil {
		return err
	}

	//// the operator will only create an ingress which acts as default ingress
	//// workaround till upstream is fixed
	ing := ingress.Config{
		// Annotations:       map[string]pulumi.StringInput{},
		ClusterIssuerType: certificate.ClusterIssuerTypeCaLocal,
		Hosts: []ingress.Host{
			{
				Name:        name,
				ServiceName: "jaeger-query",
				ServicePort: 16686,
			},
		},
		IngressClassName: ingress.IngressClassNameNginx,
		Name:             name,
		NamespaceName:    namespaceJaeger.Name,
		TLS:              true,
	}

	_, err = ingress.CreateIngress(ctx, &ing)

	if err != nil {
		return err
	}

	return nil
}

func createJaegerInstance(ctx *pulumi.Context, name string, nameSpace *namespace.Namespace) error {
	_, err := jaegerv1.NewJaeger(ctx, name, &jaegerv1.JaegerArgs{
		ApiVersion: pulumi.String("jaegertracing.io/v1"),
		Kind:       pulumi.String(strings.ToTitle(name)),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(name),
			Namespace: pulumi.String(nameSpace.Name),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

//nolint
func execGeneratedCode(ctx *pulumi.Context, namespace *namespace.Namespace) error {
	//// CHANGES: add namespace to Dep,SA
	_, err := rbacv1.NewClusterRole(ctx, "jaeger_operatorClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("jaeger-operator"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("jaegertracing.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("*"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apps"),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("jaeger-operator"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("deployments/finalizers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
					pulumi.String("persistentvolumeclaims"),
					pulumi.String("pods"),
					pulumi.String("secrets"),
					pulumi.String("serviceaccounts"),
					pulumi.String("services"),
					pulumi.String("services/finalizers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apps"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("deployments"),
					pulumi.String("daemonsets"),
					pulumi.String("replicasets"),
					pulumi.String("statefulsets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("extensions"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("batch"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("jobs"),
					pulumi.String("cronjobs"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("route.openshift.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("routes"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("console.openshift.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("consolelinks"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("autoscaling"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("horizontalpodautoscalers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("monitoring.coreos.com"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("servicemonitors"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("logging.openshift.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("elasticsearches"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("kafka.strimzi.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("kafkas"),
					pulumi.String("kafkausers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("namespaces"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apps"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("deployments"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("rbac.authorization.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("clusterrolebindings"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("patch"),
					pulumi.String("update"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "jaeger_operatorClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		Kind:       pulumi.String("ClusterRoleBinding"),
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("jaeger-operator"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind: pulumi.String("ServiceAccount"),
				Name: pulumi.String("jaeger-operator"),
				// Namespace: pulumi.String("observability"),
				Namespace: pulumi.String(namespace.Name),
			},
		},
		RoleRef: &rbacv1.RoleRefArgs{
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("jaeger-operator"),
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "jaeger_operatorDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("jaeger-operator"),
			Namespace: pulumi.String(namespace.Name),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"name": pulumi.String("jaeger-operator"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"name": pulumi.String("jaeger-operator"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: pulumi.String("jaeger-operator"),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("jaeger-operator"),
							Image: pulumi.String("jaegertracing/jaeger-operator:1.23.0"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8383),
									Name:          pulumi.String("http-metrics"),
								},
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8686),
									Name:          pulumi.String("cr-metrics"),
								},
							},
							Args: pulumi.StringArray{
								pulumi.String("start"),
							},
							ImagePullPolicy: pulumi.String("Always"),
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name:  pulumi.String("WATCH_NAMESPACE"),
									Value: pulumi.String(""),
								},
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAME"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.name"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("OPERATOR_NAME"),
									Value: pulumi.String("jaeger-operator"),
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "jaeger_operatorServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("jaeger-operator"),
			Namespace: pulumi.String(namespace.Name),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
