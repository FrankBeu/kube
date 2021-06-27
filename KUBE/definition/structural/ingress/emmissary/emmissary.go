// Package emmissary provides ambassador-crds and the default emmissary-installation
package emmissary

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1beta1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1beta1"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"thesym.site/kube/lib"
)

var (
	name               = "emmissary"
	namespaceEmmissary = &lib.Namespace{
		Name: name,
		Tier: lib.NamespaceTierEdge,
	}

	namespaceEmmissaryHosts = &lib.Namespace{
		Name: name + "-hosts",
		Tier: lib.NamespaceTierEdge,
	}
)

func CreateEmmissary(ctx *p.Context) error {

	conf := config.New(ctx, "")
	domainName := name + "." + conf.Require("domain")

	err := lib.CreateNamespaces(ctx, namespaceEmmissary, namespaceEmmissaryHosts)
	if err != nil {
		return err
	}

	// TODO: warning: apiextensions.k8s.io/v1beta1/CustomResourceDefinition is deprecated by
	// apiextensions.k8s.io/v1/CustomResourceDefinition and not supported by Kubernetes v1.22+ clusters.
	// has to be solved upstream - warnings are supprest
	emmissaryCrds := []lib.Crd{
		{Name: "emmissary-ambassador-definition", Location: "./crds/emmissary/cdrDefinitions/ambassador-crds.yaml"},
	}

	err = lib.RegisterCRDs(ctx, emmissaryCrds)
	if err != nil {
		return err
	}

	// createResourcesForDiagnostics(ctx, domainName)

	err = execGeneratedCode(ctx)
	if err != nil {
		return err
	}

	return nil

}

// createResourcesForDiagnostics provides all resources for the emmissaryDiagnosticsFrontend

//  func createResourcesForDiagnostics(ctx *p.Context, domainName string) error {
//
//  	_, err := getambassadorv2.NewHost(ctx, "emmissaryDiagnostic", &getambassadorv2.HostArgs{
//  		// ApiVersion: p.String("getambassador.io/v2"),
//  		// Kind:       p.String("Host"),
//  		// // Metadata:   &getambassadorv2.HostMetadataArgs{},
//  		// Metadata: &metav1.ObjectMetaArgs{
//  		// 	Name:      p.String("emmissaryDiagnostic"),
//  		// 	Namespace: p.String(namespaceEmmissaryHosts.Name),
//  		// },
//
//  		// Spec: &getambassadorv2.HostSpecArgs{
//  		// //	Hostname: p.String("emdia."stage.thesym.site"),
//  		// 	Hostname: p.String(domainName),
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

func execGeneratedCode(ctx *p.Context) error {
	_, err := corev1.NewService(ctx, "defaultAmbassador_adminService", &corev1.ServiceArgs{
		ApiVersion: p.String("v1"),
		Kind:       p.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      p.String("ambassador-admin"),
			Namespace: p.String(namespaceEmmissary.Name),
			Labels: p.StringMap{
				"service": p.String("ambassador-admin"),
				"product": p.String("aes"),
			},
			Annotations: p.StringMap{
				"a8r.io/owner":         p.String("Ambassador Labs"),
				"a8r.io/repository":    p.String("github.com/datawire/ambassador"),
				"a8r.io/description":   p.String("The Ambassador Edge Stack admin service for internal use and health checks."),
				"a8r.io/documentation": p.String("https://www.getambassador.io/docs/edge-stack/latest/"),
				"a8r.io/chat":          p.String("http://a8r.io/Slack"),
				"a8r.io/bugs":          p.String("https://github.com/datawire/ambassador/issues"),
				"a8r.io/support":       p.String("https://www.getambassador.io/about-us/support/"),
				"a8r.io/dependencies":  p.String("None"),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: p.String("NodePort"),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       p.Int(8877),
					TargetPort: p.String("admin"),
					Protocol:   p.String("TCP"),
					Name:       p.String("ambassador-admin"),
				},
				&corev1.ServicePortArgs{
					Port:       p.Int(8005),
					TargetPort: p.Int(8005),
					Protocol:   p.String("TCP"),
					Name:       p.String("ambassador-snapshot"),
				},
			},
			Selector: p.StringMap{
				"service": p.String("ambassador"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassadorClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador"),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
		AggregationRule: &rbacv1beta1.AggregationRuleArgs{
			ClusterRoleSelectors: metav1.LabelSelectorArray{
				&metav1.LabelSelectorArgs{
					MatchLabels: p.StringMap{
						"rbac.getambassador.io/role-group": p.String("ambassador"),
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
		ApiVersion: p.String("v1"),
		Kind:       p.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      p.String("ambassador"),
			Namespace: p.String(namespaceEmmissary.Name),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRoleBinding(ctx, "ambassadorClusterRoleBinding", &rbacv1beta1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador"),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
		RoleRef: &rbacv1beta1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("ambassador"),
		},
		Subjects: rbacv1beta1.SubjectArray{
			&rbacv1beta1.SubjectArgs{
				Name:      p.String("ambassador"),
				Namespace: p.String(namespaceEmmissary.Name),
				Kind:      p.String("ServiceAccount"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_crdClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador-crd"),
			Labels: p.StringMap{
				"product":                          p.String("aes"),
				"rbac.getambassador.io/role-group": p.String("ambassador"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("apiextensions.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("customresourcedefinitions"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("delete"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_watchClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador-watch"),
			Labels: p.StringMap{
				"product":                          p.String("aes"),
				"rbac.getambassador.io/role-group": p.String("ambassador"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("namespaces"),
					p.String("services"),
					p.String("secrets"),
					p.String("endpoints"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("getambassador.io"),
				},
				Resources: p.StringArray{
					p.String("*"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("update"),
					p.String("patch"),
					p.String("create"),
					p.String("delete"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("getambassador.io"),
				},
				Resources: p.StringArray{
					p.String("mappings/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("networking.internal.knative.dev"),
				},
				Resources: p.StringArray{
					p.String("clusteringresses"),
					p.String("ingresses"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("networking.x-k8s.io"),
				},
				Resources: p.StringArray{
					p.String("*"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("networking.internal.knative.dev"),
				},
				Resources: p.StringArray{
					p.String("ingresses/status"),
					p.String("clusteringresses/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("extensions"),
					p.String("networking.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("ingresses"),
					p.String("ingressclasses"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("extensions"),
					p.String("networking.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("ingresses/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "defaultAmbassadorDeployment", &appsv1.DeploymentArgs{
		ApiVersion: p.String("apps/v1"),
		Kind:       p.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      p.String("ambassador"),
			Namespace: p.String(namespaceEmmissary.Name),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: p.Int(3),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: p.StringMap{
					"service": p.String("ambassador"),
				},
			},
			Strategy: &appsv1.DeploymentStrategyArgs{
				Type: p.String("RollingUpdate"),
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: p.StringMap{
						"service":                      p.String("ambassador"),
						"app.kubernetes.io/managed-by": p.String("getambassador.io"),
					},
					Annotations: p.StringMap{
						"consul.hashicorp.com/connect-inject": p.String("false"),
						"sidecar.istio.io/inject":             p.String("false"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					TerminationGracePeriodSeconds: p.Int(0),
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsUser: p.Int(8888),
					},
					RestartPolicy:      p.String("Always"),
					ServiceAccountName: p.String("ambassador"),
					Volumes: corev1.VolumeArray{
						&corev1.VolumeArgs{
							Name: p.String("ambassador-pod-info"),
							DownwardAPI: &corev1.DownwardAPIVolumeSourceArgs{
								Items: corev1.DownwardAPIVolumeFileArray{
									&corev1.DownwardAPIVolumeFileArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: p.String("metadata.labels"),
										},
										Path: p.String("labels"),
									},
								},
							},
						},
					},
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            p.String("ambassador"),
							Image:           p.String("docker.io/datawire/ambassador:1.13.8"),
							ImagePullPolicy: p.String("IfNotPresent"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									Name:          p.String("http"),
									ContainerPort: p.Int(8080),
								},
								&corev1.ContainerPortArgs{
									Name:          p.String("https"),
									ContainerPort: p.Int(8443),
								},
								&corev1.ContainerPortArgs{
									Name:          p.String("admin"),
									ContainerPort: p.Int(8877),
								},
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: p.String("HOST_IP"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: p.String("status.hostIP"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name: p.String("AMBASSADOR_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: p.String("metadata.namespace"),
										},
									},
								},
							},
							SecurityContext: &corev1.SecurityContextArgs{
								AllowPrivilegeEscalation: p.Bool(false),
							},
							LivenessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: p.String("/ambassador/v0/check_alive"),
									Port: p.String("admin"),
								},
								FailureThreshold:    p.Int(3),
								InitialDelaySeconds: p.Int(30),
								PeriodSeconds:       p.Int(3),
							},
							ReadinessProbe: &corev1.ProbeArgs{
								HttpGet: &corev1.HTTPGetActionArgs{
									Path: p.String("/ambassador/v0/check_ready"),
									Port: p.String("admin"),
								},
								FailureThreshold:    p.Int(3),
								InitialDelaySeconds: p.Int(30),
								PeriodSeconds:       p.Int(3),
							},
							VolumeMounts: corev1.VolumeMountArray{
								&corev1.VolumeMountArgs{
									Name:      p.String("ambassador-pod-info"),
									MountPath: p.String("/tmp/ambassador-pod-info"),
									ReadOnly:  p.Bool(true),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Limits: p.StringMap{
									"cpu":    p.String("1"),
									"memory": p.String("400Mi"),
								},
								Requests: p.StringMap{
									"cpu":    p.String("200m"),
									"memory": p.String("100Mi"),
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
											MatchLabels: p.StringMap{
												"service": p.String("ambassador"),
											},
										},
										TopologyKey: p.String("kubernetes.io/hostname"),
									},
									Weight: p.Int(100),
								},
							},
						},
					},
					ImagePullSecrets: corev1.LocalObjectReferenceArray{},
					DnsPolicy:        p.String("ClusterFirst"),
					HostNetwork:      p.Bool(false),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "defaultAmbassador_agentServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: p.String("v1"),
		Kind:       p.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      p.String("ambassador-agent"),
			Namespace: p.String(namespaceEmmissary.Name),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRoleBinding(ctx, "ambassador_agentClusterRoleBinding", &rbacv1beta1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador-agent"),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
		RoleRef: &rbacv1beta1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("ambassador-agent"),
		},
		Subjects: rbacv1beta1.SubjectArray{
			&rbacv1beta1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("ambassador-agent"),
				Namespace: p.String("default"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_agentClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador-agent"),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
		AggregationRule: &rbacv1beta1.AggregationRuleArgs{
			ClusterRoleSelectors: metav1.LabelSelectorArray{
				&metav1.LabelSelectorArgs{
					MatchLabels: p.StringMap{
						"rbac.getambassador.io/role-group": p.String("ambassador-agent"),
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
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador-agent-pods"),
			Labels: p.StringMap{
				"rbac.getambassador.io/role-group": p.String("ambassador-agent"),
				"product":                          p.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("pods"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_agent_rolloutsClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador-agent-rollouts"),
			Labels: p.StringMap{
				"rbac.getambassador.io/role-group": p.String("ambassador-agent"),
				"product":                          p.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("argoproj.io"),
				},
				Resources: p.StringArray{
					p.String("rollouts"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewClusterRole(ctx, "ambassador_agent_applicationsClusterRole", &rbacv1beta1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("ambassador-agent-applications"),
			Labels: p.StringMap{
				"rbac.getambassador.io/role-group": p.String("ambassador-agent"),
				"product":                          p.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("argoproj.io"),
				},
				Resources: p.StringArray{
					p.String("applications"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewRole(ctx, "defaultAmbassador_agent_configRole", &rbacv1beta1.RoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      p.String("ambassador-agent-config"),
			Namespace: p.String(namespaceEmmissary.Name),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
		Rules: rbacv1beta1.PolicyRuleArray{
			&rbacv1beta1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("configmaps"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1beta1.NewRoleBinding(ctx, "defaultAmbassador_agent_configRoleBinding", &rbacv1beta1.RoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1beta1"),
		Kind:       p.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      p.String("ambassador-agent-config"),
			Namespace: p.String(namespaceEmmissary.Name),
			Labels: p.StringMap{
				"product": p.String("aes"),
			},
		},
		RoleRef: &rbacv1beta1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("Role"),
			Name:     p.String("ambassador-agent-config"),
		},
		Subjects: rbacv1beta1.SubjectArray{
			&rbacv1beta1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("ambassador-agent"),
				Namespace: p.String("default"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "defaultAmbassador_agentDeployment", &appsv1.DeploymentArgs{
		ApiVersion: p.String("apps/v1"),
		Kind:       p.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      p.String("ambassador-agent"),
			Namespace: p.String(namespaceEmmissary.Name),
			Labels: p.StringMap{
				"app.kubernetes.io/name":     p.String("ambassador-agent"),
				"app.kubernetes.io/instance": p.String("ambassador"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: p.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: p.StringMap{
					"app.kubernetes.io/name":     p.String("ambassador-agent"),
					"app.kubernetes.io/instance": p.String("ambassador"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: p.StringMap{
						"app.kubernetes.io/name":     p.String("ambassador-agent"),
						"app.kubernetes.io/instance": p.String("ambassador"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: p.String("ambassador-agent"),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            p.String("agent"),
							Image:           p.String("docker.io/datawire/ambassador:1.13.8"),
							ImagePullPolicy: p.String("IfNotPresent"),
							Command: p.StringArray{
								p.String("agent"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: p.String("AGENT_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: p.String("metadata.namespace"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name:  p.String("AGENT_CONFIG_RESOURCE_NAME"),
									Value: p.String("ambassador-agent-cloud-token"),
								},
								&corev1.EnvVarArgs{
									Name:  p.String("RPC_CONNECTION_ADDRESS"),
									Value: p.String("https://app.getambassador.io/"),
								},
								&corev1.EnvVarArgs{
									Name:  p.String("AES_SNAPSHOT_URL"),
									Value: p.String("http://ambassador-admin.default:8005/snapshot-external"),
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
