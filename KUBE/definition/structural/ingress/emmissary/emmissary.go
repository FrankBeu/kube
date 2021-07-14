// Package emmissary provides ambassador-crds and the default emmissary-installation
package emmissary

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1beta1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1beta1"
	pulumi "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/crd"
	"thesym.site/kube/lib/namespace"
)

var (
	name               = "emmissary"
	namespaceEmmissary = &namespace.Namespace{
		Name: name,
		Tier: namespace.NamespaceTierEdge,
	}

	namespaceEmmissaryHosts = &namespace.Namespace{
		Name: name + "-hosts",
		Tier: namespace.NamespaceTierEdge,
	}
)

func CreateEmmissary(ctx *pulumi.Context) error {
	_, err := namespace.CreateNamespace(ctx, namespaceEmmissary)
	if err != nil {
		return err
	}

	_, err = namespace.CreateNamespace(ctx, namespaceEmmissaryHosts)
	if err != nil {
		return err
	}

	// TODO: warning: apiextensions.k8s.io/v1beta1/CustomResourceDefinition is deprecated by
	// apiextensions.k8s.io/v1/CustomResourceDefinition and not supported by Kubernetes v1.22+ clusters.
	// has to be solved upstream - warnings are supprest
	emmissaryCrds := []crd.Crd{
		{Name: "emmissary-ambassador-definition", Location: "./crds/emmissary/cdrDefinitions/ambassador-crds.yaml"},
	}

	err = crd.RegisterCRDs(ctx, emmissaryCrds)
	if err != nil {
		return err
	}

	//nolint:gocritic
	// createResourcesForDiagnostics(ctx, domainName)

	err = execGeneratedCode(ctx)
	if err != nil {
		return err
	}

	return nil
}

// createResourcesForDiagnostics provides all resources for the emmissaryDiagnosticsFrontend

//  func createResourcesForDiagnostics(ctx *pulumi.Context, domainName string) error {
//
//  	_, err := getambassadorv2.NewHost(ctx, "emmissaryDiagnostic", &getambassadorv2.HostArgs{
//  		// ApiVersion: pulumi.String("getambassador.io/v2"),
//  		// Kind:       pulumi.String("Host"),
//  		// // Metadata:   &getambassadorv2.HostMetadataArgs{},
//  		// Metadata: &metav1.ObjectMetaArgs{
//  		// 	Name:      pulumi.String("emmissaryDiagnostic"),
//  		// 	Namespace: pulumi.String(namespaceEmmissaryHosts.Name),
//  		// },
//
//  		// Spec: &getambassadorv2.HostSpecArgs{
//  		// //	Hostname: pulumi.String("emdia."stage.thesym.site"),
//  		// 	Hostname: pulumi.String(domainName),
//  		// },
//  		// Status: nil,
//  	})
//
//  	if err != nil {
//  		return err
//  	}
//
//  	////////////////////////////////////////////////////////////////////////
//  	// apiVersion: getambassador.io/v2
//  	// kind: Mapping
//  	// metadata:
//  	//   name: gitea
//  	//   namespace: gitea
//  	// spec:
//  	//   prefix: /
//  	//   #service: gitea-service.gitea:3000
//  	//   service: gitea-http.gitea:3000
//  	//   host: gitea.thesym.site
//  	//   timeout_ms: 60000// apiVersion: getambassador.io/v2
//  	// ---
//  	// kind: Host
//  	// metadata:
//  	//   name: gitea
//  	//   namespace: ambassador-hosts
//  	// spec:
//  	//   hostname: gitea.thesym.site
//  	//   acmeProvider:
//  	//     authority: https://acme-v02.api.letsencrypt.org/directory
//  	//     email: fbeutelschiess@gmail.com
//  	//     privateKeySecret:
//  	//       name: gitea-key
//  	//   tlsSecret:
//  	//     name: gitea
//  	//   requestPolicy:
//  	//     insecure://       action: Redirect
//
//  	return nil
//  }

//nolint // generatedCode
func execGeneratedCode(ctx *pulumi.Context) error {
	_, err := corev1.NewService(ctx, "defaultAmbassador_adminService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ambassador-admin"),
			Namespace: pulumi.String(namespaceEmmissary.Name),
			Labels: pulumi.StringMap{
				"service": pulumi.String("ambassador-admin"),
				"product": pulumi.String("aes"),
			},
			Annotations: pulumi.StringMap{
				"a8r.io/owner":         pulumi.String("Ambassador Labs"),
				"a8r.io/repository":    pulumi.String("github.com/datawire/ambassador"),
				"a8r.io/description":   pulumi.String("The Ambassador Edge Stack admin service for internal use and health checks."),
				"a8r.io/documentation": pulumi.String("https://www.getambassador.io/docs/edge-stack/latest/"),
				"a8r.io/chat":          pulumi.String("http://a8r.io/Slack"),
				"a8r.io/bugs":          pulumi.String("https://github.com/datawire/ambassador/issues"),
				"a8r.io/support":       pulumi.String("https://www.getambassador.io/about-us/support/"),
				"a8r.io/dependencies":  pulumi.String("None"),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: pulumi.String("NodePort"),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(8877),
					TargetPort: pulumi.String("admin"),
					Protocol:   pulumi.String("TCP"),
					Name:       pulumi.String("ambassador-admin"),
				},
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(8005),
					TargetPort: pulumi.Int(8005),
					Protocol:   pulumi.String("TCP"),
					Name:       pulumi.String("ambassador-snapshot"),
				},
			},
			Selector: pulumi.StringMap{
				"service": pulumi.String("ambassador"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassadorClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador"),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
		AggregationRule: &rbacv1beta1.AggregationRuleArgs{
			ClusterRoleSelectors: metav1.LabelSelectorArray{
				&metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"rbac.getambassador.io/role-group": pulumi.String("ambassador"),
					},
				},
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "defaultAmbassadorServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ambassador"),
			Namespace: pulumi.String(namespaceEmmissary.Name),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRoleBinding(ctx, "ambassadorClusterRoleBinding", &rbacv1beta1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador"),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
		RoleRef: &rbacv1beta1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("ambassador"),
		},
		Subjects: rbacv1beta1.SubjectArray{
			&rbacv1beta1.SubjectArgs{
				Name:      pulumi.String("ambassador"),
				Namespace: pulumi.String(namespaceEmmissary.Name),
				Kind:      pulumi.String("ServiceAccount"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_crdClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador-crd"),
			Labels: pulumi.StringMap{
				"product":                          pulumi.String("aes"),
				"rbac.getambassador.io/role-group": pulumi.String("ambassador"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apiextensions.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("customresourcedefinitions"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("delete"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_watchClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador-watch"),
			Labels: pulumi.StringMap{
				"product":                          pulumi.String("aes"),
				"rbac.getambassador.io/role-group": pulumi.String("ambassador"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("namespaces"),
					pulumi.String("services"),
					pulumi.String("secrets"),
					pulumi.String("endpoints"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("getambassador.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("*"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("update"),
					pulumi.String("patch"),
					pulumi.String("create"),
					pulumi.String("delete"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("getambassador.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("mappings/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.internal.knative.dev"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("clusteringresses"),
					pulumi.String("ingresses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.x-k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("*"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.internal.knative.dev"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses/status"),
					pulumi.String("clusteringresses/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("extensions"),
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses"),
					pulumi.String("ingressclasses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("extensions"),
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "defaultAmbassadorDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ambassador"),
			Namespace: pulumi.String(namespaceEmmissary.Name),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(3),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"service": pulumi.String("ambassador"),
				},
			},
			Strategy: &appsv1.DeploymentStrategyArgs{
				Type: pulumi.String("RollingUpdate"),
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"service":                      pulumi.String("ambassador"),
						"app.kubernetes.io/managed-by": pulumi.String("getambassador.io"),
					},
					Annotations: pulumi.StringMap{
						"consul.hashicorp.com/connect-inject": pulumi.String("false"),
						"sidecar.istio.io/inject":             pulumi.String("false"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					TerminationGracePeriodSeconds: pulumi.Int(0),
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsUser: pulumi.Int(8888),
					},
					RestartPolicy:      pulumi.String("Always"),
					ServiceAccountName: pulumi.String("ambassador"),
					Volumes: corev1.VolumeArray{
						&corev1.VolumeArgs{
							Name: pulumi.String("ambassador-pod-info"),
							DownwardAPI: &corev1.DownwardAPIVolumeSourceArgs{
								Items: corev1.DownwardAPIVolumeFileArray{
									&corev1.DownwardAPIVolumeFileArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.labels"),
										},
										Path: pulumi.String("labels"),
									},
								},
							},
						},
					},
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String("ambassador"),
							Image:           pulumi.String("docker.io/datawire/ambassador:1.13.8"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("http"),
									ContainerPort: pulumi.Int(8080),
								},
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("https"),
									ContainerPort: pulumi.Int(8443),
								},
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("admin"),
									ContainerPort: pulumi.Int(8877),
								},
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("HOST_IP"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("status.hostIP"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name: pulumi.String("AMBASSADOR_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
							},
							SecurityContext: &corev1.SecurityContextArgs{
								AllowPrivilegeEscalation: pulumi.Bool(false),
							},
							LivenessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/ambassador/v0/check_alive"),
									Port: pulumi.String("admin"),
								},
								FailureThreshold:    pulumi.Int(3),
								InitialDelaySeconds: pulumi.Int(30),
								PeriodSeconds:       pulumi.Int(3),
							},
							ReadinessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: pulumi.String("/ambassador/v0/check_ready"),
									Port: pulumi.String("admin"),
								},
								FailureThreshold:    pulumi.Int(3),
								InitialDelaySeconds: pulumi.Int(30),
								PeriodSeconds:       pulumi.Int(3),
							},
							VolumeMounts: corev1.VolumeMountArray{
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("ambassador-pod-info"),
									MountPath: pulumi.String("/tmp/ambassador-pod-info"),
									ReadOnly:  pulumi.Bool(true),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Limits: pulumi.StringMap{
									"cpu":    pulumi.String("1"),
									"memory": pulumi.String("400Mi"),
								},
								Requests: pulumi.StringMap{
									"cpu":    pulumi.String("200m"),
									"memory": pulumi.String("100Mi"),
								},
							},
						},
					},
					Affinity: &corev1.AffinityArgs{
						PodAntiAffinity: &corev1.PodAntiAffinityArgs{
							PreferredDuringSchedulingIgnoredDuringExecution: corev1.WeightedPodAffinityTermArray{
								&corev1.WeightedPodAffinityTermArgs{
									PodAffinityTerm: &corev1.PodAffinityTermArgs{
										LabelSelector: &metav1.LabelSelectorArgs{
											MatchLabels: pulumi.StringMap{
												"service": pulumi.String("ambassador"),
											},
										},
										TopologyKey: pulumi.String("kubernetes.io/hostname"),
									},
									Weight: pulumi.Int(100),
								},
							},
						},
					},
					ImagePullSecrets: corev1.LocalObjectReferenceArray{},
					DnsPolicy:        pulumi.String("ClusterFirst"),
					HostNetwork:      pulumi.Bool(false),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "defaultAmbassador_agentServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ambassador-agent"),
			Namespace: pulumi.String(namespaceEmmissary.Name),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRoleBinding(ctx, "ambassador_agentClusterRoleBinding", &rbacv1beta1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador-agent"),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
		RoleRef: &rbacv1beta1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("ambassador-agent"),
		},
		Subjects: rbacv1beta1.SubjectArray{
			&rbacv1beta1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("ambassador-agent"),
				Namespace: pulumi.String("default"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_agentClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador-agent"),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
		AggregationRule: &rbacv1beta1.AggregationRuleArgs{
			ClusterRoleSelectors: metav1.LabelSelectorArray{
				&metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"rbac.getambassador.io/role-group": pulumi.String("ambassador-agent"),
					},
				},
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_agent_podsClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador-agent-pods"),
			Labels: pulumi.StringMap{
				"rbac.getambassador.io/role-group": pulumi.String("ambassador-agent"),
				"product":                          pulumi.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("pods"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_agent_rolloutsClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador-agent-rollouts"),
			Labels: pulumi.StringMap{
				"rbac.getambassador.io/role-group": pulumi.String("ambassador-agent"),
				"product":                          pulumi.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("argoproj.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("rollouts"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_agent_applicationsClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ambassador-agent-applications"),
			Labels: pulumi.StringMap{
				"rbac.getambassador.io/role-group": pulumi.String("ambassador-agent"),
				"product":                          pulumi.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("argoproj.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("applications"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewRole(ctx, "defaultAmbassador_agent_configRole", &rbacv1beta1.RoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ambassador-agent-config"),
			Namespace: pulumi.String(namespaceEmmissary.Name),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewRoleBinding(ctx, "defaultAmbassador_agent_configRoleBinding", &rbacv1beta1.RoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       pulumi.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ambassador-agent-config"),
			Namespace: pulumi.String(namespaceEmmissary.Name),
			Labels: pulumi.StringMap{
				"product": pulumi.String("aes"),
			},
		},
		RoleRef: &rbacv1beta1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("Role"),
			Name:     pulumi.String("ambassador-agent-config"),
		},
		Subjects: rbacv1beta1.SubjectArray{
			&rbacv1beta1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("ambassador-agent"),
				Namespace: pulumi.String("default"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "defaultAmbassador_agentDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("ambassador-agent"),
			Namespace: pulumi.String(namespaceEmmissary.Name),
			Labels: pulumi.StringMap{
				"app.kubernetes.io/name":     pulumi.String("ambassador-agent"),
				"app.kubernetes.io/instance": pulumi.String("ambassador"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app.kubernetes.io/name":     pulumi.String("ambassador-agent"),
					"app.kubernetes.io/instance": pulumi.String("ambassador"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app.kubernetes.io/name":     pulumi.String("ambassador-agent"),
						"app.kubernetes.io/instance": pulumi.String("ambassador"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: pulumi.String("ambassador-agent"),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String("agent"),
							Image:           pulumi.String("docker.io/datawire/ambassador:1.13.8"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Command: pulumi.StringArray{
								pulumi.String("agent"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("AGENT_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("AGENT_CONFIG_RESOURCE_NAME"),
									Value: pulumi.String("ambassador-agent-cloud-token"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("RPC_CONNECTION_ADDRESS"),
									Value: pulumi.String("https://app.getambassador.io/"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("AES_SNAPSHOT_URL"),
									Value: pulumi.String("http://ambassador-admin.default:8005/snapshot-external"),
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
	return nil
}
