// Package nginx provides the default nginxIngressController-installation
package nginx

import (
	"fmt"

	admissionregistrationv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/admissionregistration/v1"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/batch/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib"
)

var (
	ingressControllerPortHTTP  = 30080
	ingressControllerPortHTTPS = 30443
	namespaceNginxIngress      = &lib.Namespace{
		Name: "ingress-nginx",
		Tier: lib.NamespaceTierEdge,
		AdditionalLabels: []lib.NamespaceLabel{
			{Name: "app.kubernetes.io/name", Value: "ingress-nginx"},
			{Name: "app.kubernetes.io/instance", Value: "ingress-nginx"},
		},
	}
)

func CreateNginxIngressController(ctx *pulumi.Context) error {
	err := lib.CreateNamespaces(ctx, namespaceNginxIngress)
	if err != nil {
		return err
	}

	err = execGeneratedCode(ctx)
	if err != nil {
		return err
	}

	return nil
}

//nolint
func execGeneratedCode(ctx *pulumi.Context) error {
	// CHANGES: namespaceChange Resource Name; DefaultLabels

	_, err := rbacv1.NewClusterRole(ctx, "ingress_nginxClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
			},
			Name: pulumi.String("ingress-nginx"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
					pulumi.String("endpoints"),
					pulumi.String("nodes"),
					pulumi.String("pods"),
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("nodes"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("services"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("extensions"),
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
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
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingressclasses"),
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
	_, err = rbacv1.NewClusterRoleBinding(ctx, "ingress_nginxClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
			},
			Name: pulumi.String("ingress-nginx"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("ingress-nginx"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("ingress-nginx"),
				Namespace: pulumi.String("ingress-nginx"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewConfigMap(ctx, "ingress_nginxIngress_nginx_controllerConfigMap", &corev1.ConfigMapArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ConfigMap"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
			Name:      pulumi.String("ingress-nginx-controller"),
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		Data: nil,
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "ingress_nginxIngress_nginx_controllerDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
			Name:      pulumi.String("ingress-nginx-controller"),
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app.kubernetes.io/name":      pulumi.String("ingress-nginx"),
					"app.kubernetes.io/instance":  pulumi.String("ingress-nginx"),
					"app.kubernetes.io/component": pulumi.String("controller"),
				},
			},
			RevisionHistoryLimit: pulumi.Int(10),
			MinReadySeconds:      pulumi.Int(0),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app.kubernetes.io/name":      pulumi.String("ingress-nginx"),
						"app.kubernetes.io/instance":  pulumi.String("ingress-nginx"),
						"app.kubernetes.io/component": pulumi.String("controller"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					DnsPolicy: pulumi.String("ClusterFirst"),
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String("controller"),
							Image:           pulumi.String("k8s.gcr.io/ingress-nginx/controller:v0.46.0@sha256:52f0058bed0a17ab0fb35628ba97e8d52b5d32299fbc03cc0f6c7b9ff036b61a"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Lifecycle: &corev1.LifecycleArgs{
								PreStop: &corev1.HandlerArgs{
									Exec: &corev1.ExecActionArgs{
										Command: pulumi.StringArray{
											pulumi.String("/wait-shutdown"),
										},
									},
								},
							},
							Args: pulumi.StringArray{
								pulumi.String("/nginx-ingress-controller"),
								pulumi.String("--election-id=ingress-controller-leader"),
								pulumi.String("--ingress-class=nginx"),
								pulumi.String(fmt.Sprintf("%v%v%v", "--configmap=", "$", "(POD_NAMESPACE)/ingress-nginx-controller")),
								pulumi.String("--validating-webhook=:8443"),
								pulumi.String("--validating-webhook-certificate=/usr/local/certificates/cert"),
								pulumi.String("--validating-webhook-key=/usr/local/certificates/key"),
							},
							SecurityContext: &corev1.SecurityContextArgs{
								Capabilities: &corev1.CapabilitiesArgs{
									Drop: pulumi.StringArray{
										pulumi.String("ALL"),
									},
									Add: pulumi.StringArray{
										pulumi.String("NET_BIND_SERVICE"),
									},
								},
								RunAsUser:                pulumi.Int(101),
								AllowPrivilegeEscalation: pulumi.Bool(true),
							},
							Env: corev1.EnvVarArray{
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
									Name:  pulumi.String("LD_PRELOAD"),
									Value: pulumi.String("/usr/local/lib/libmimalloc.so"),
								},
							},
							LivenessProbe: &corev1.ProbeArgs{
								FailureThreshold: pulumi.Int(5),
								HttpGet: &corev1.HTTPGetActionArgs{
									Path:   pulumi.String("/healthz"),
									Port:   pulumi.Int(10254),
									Scheme: pulumi.String("HTTP"),
								},
								InitialDelaySeconds: pulumi.Int(10),
								PeriodSeconds:       pulumi.Int(10),
								SuccessThreshold:    pulumi.Int(1),
								TimeoutSeconds:      pulumi.Int(1),
							},
							ReadinessProbe: &corev1.ProbeArgs{
								FailureThreshold: pulumi.Int(3),
								HttpGet: &corev1.HTTPGetActionArgs{
									Path:   pulumi.String("/healthz"),
									Port:   pulumi.Int(10254),
									Scheme: pulumi.String("HTTP"),
								},
								InitialDelaySeconds: pulumi.Int(10),
								PeriodSeconds:       pulumi.Int(10),
								SuccessThreshold:    pulumi.Int(1),
								TimeoutSeconds:      pulumi.Int(1),
							},
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("http"),
									ContainerPort: pulumi.Int(80),
									Protocol:      pulumi.String("TCP"),
								},
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("https"),
									ContainerPort: pulumi.Int(443),
									Protocol:      pulumi.String("TCP"),
								},
								&corev1.ContainerPortArgs{
									Name:          pulumi.String("webhook"),
									ContainerPort: pulumi.Int(8443),
									Protocol:      pulumi.String("TCP"),
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("webhook-cert"),
									MountPath: pulumi.String("/usr/local/certificates/"),
									ReadOnly:  pulumi.Bool(true),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"cpu":    pulumi.String("100m"),
									"memory": pulumi.String("90Mi"),
								},
							},
						},
					},
					NodeSelector: pulumi.StringMap{
						"kubernetes.io/os": pulumi.String("linux"),
					},
					ServiceAccountName:            pulumi.String("ingress-nginx"),
					TerminationGracePeriodSeconds: pulumi.Int(300),
					Volumes: corev1.VolumeArray{
						&corev1.VolumeArgs{
							Name: pulumi.String("webhook-cert"),
							Secret: &corev1.SecretVolumeSourceArgs{
								SecretName: pulumi.String("ingress-nginx-admission"),
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
	_, err = rbacv1.NewRole(ctx, "ingress_nginxIngress_nginxRole", &rbacv1.RoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
			Name:      pulumi.String("ingress-nginx"),
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("namespaces"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
					pulumi.String("pods"),
					pulumi.String("secrets"),
					pulumi.String("endpoints"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("services"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("extensions"),
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
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
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingressclasses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("ingress-controller-leader-nginx"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "ingress_nginxIngress_nginxRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
			Name:      pulumi.String("ingress-nginx"),
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("Role"),
			Name:     pulumi.String("ingress-nginx"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("ingress-nginx"),
				Namespace: pulumi.String(namespaceNginxIngress.Name),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "ingress_nginxIngress_nginx_controller_admissionService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
			Name:      pulumi.String("ingress-nginx-controller-admission"),
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: pulumi.String("ClusterIP"),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("https-webhook"),
					Port:       pulumi.Int(443),
					TargetPort: pulumi.String("webhook"),
				},
			},
			Selector: pulumi.StringMap{
				"app.kubernetes.io/name":      pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance":  pulumi.String("ingress-nginx"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "ingress_nginxIngress_nginx_controllerService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Annotations: nil,
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
			Name:      pulumi.String("ingress-nginx-controller"),
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: pulumi.String("NodePort"),
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("http"),
					NodePort:   pulumi.Int(ingressControllerPortHTTP),
					Port:       pulumi.Int(80),
					Protocol:   pulumi.String("TCP"),
					TargetPort: pulumi.String("http"),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("https"),
					NodePort:   pulumi.Int(ingressControllerPortHTTPS),
					Port:       pulumi.Int(443),
					Protocol:   pulumi.String("TCP"),
					TargetPort: pulumi.String("https"),
				},
			},
			Selector: pulumi.StringMap{
				"app.kubernetes.io/name":      pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance":  pulumi.String("ingress-nginx"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "ingress_nginxIngress_nginxServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("controller"),
			},
			Name:      pulumi.String("ingress-nginx"),
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		AutomountServiceAccountToken: pulumi.Bool(true),
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "ingress_nginxIngress_nginx_admissionServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ingress-nginx-admission"),
			// Annotations: pulumi.StringMap{
			// "helm.sh/hook":               pulumi.String("pre-install,pre-upgrade,post-install,post-upgrade"),
			// "helm.sh/hook-delete-policy": pulumi.String("before-hook-creation,hook-succeeded"),
			// },
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
	})
	if err != nil {
		return err
	}
	_, err = admissionregistrationv1.NewValidatingWebhookConfiguration(ctx, "ingress_nginx_admissionValidatingWebhookConfiguration", &admissionregistrationv1.ValidatingWebhookConfigurationArgs{
		ApiVersion: pulumi.String("admissionregistration.k8s.io/v1"),
		Kind:       pulumi.String("ValidatingWebhookConfiguration"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
			Name: pulumi.String("ingress-nginx-admission"),
		},
		Webhooks: admissionregistrationv1.ValidatingWebhookArray{
			&admissionregistrationv1.ValidatingWebhookArgs{
				Name:        pulumi.String("validate.nginx.ingress.kubernetes.io"),
				MatchPolicy: pulumi.String("Equivalent"),
				Rules: admissionregistrationv1.RuleWithOperationsArray{
					&admissionregistrationv1.RuleWithOperationsArgs{
						ApiGroups: pulumi.StringArray{
							pulumi.String("networking.k8s.io"),
						},
						ApiVersions: pulumi.StringArray{
							pulumi.String("v1beta1"),
						},
						Operations: pulumi.StringArray{
							pulumi.String("CREATE"),
							pulumi.String("UPDATE"),
						},
						Resources: pulumi.StringArray{
							pulumi.String("ingresses"),
						},
					},
				},
				FailurePolicy: pulumi.String("Fail"),
				SideEffects:   pulumi.String("None"),
				AdmissionReviewVersions: pulumi.StringArray{
					pulumi.String("v1"),
					pulumi.String("v1beta1"),
				},
				ClientConfig: &admissionregistrationv1.WebhookClientConfigArgs{
					Service: &admissionregistrationv1.ServiceReferenceArgs{
						Namespace: pulumi.String(namespaceNginxIngress.Name),
						Name:      pulumi.String("ingress-nginx-controller-admission"),
						Path:      pulumi.String("/networking/v1beta1/ingresses"),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "ingress_nginx_admissionClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ingress-nginx-admission"),
			// Annotations: pulumi.StringMap{
			// "helm.sh/hook":               pulumi.String("pre-install,pre-upgrade,post-install,post-upgrade"),
			// "helm.sh/hook-delete-policy": pulumi.String("before-hook-creation,hook-succeeded"),
			// },
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("admissionregistration.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("validatingwebhookconfigurations"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("update"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "ingress_nginx_admissionClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ingress-nginx-admission"),
			// Annotations: pulumi.StringMap{
			// "helm.sh/hook":               pulumi.String("pre-install,pre-upgrade,post-install,post-upgrade"),
			// "helm.sh/hook-delete-policy": pulumi.String("before-hook-creation,hook-succeeded"),
			// },
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("ingress-nginx-admission"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("ingress-nginx-admission"),
				Namespace: pulumi.String(namespaceNginxIngress.Name),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = batchv1.NewJob(ctx, "ingress_nginxIngress_nginx_admission_createJob", &batchv1.JobArgs{
		ApiVersion: pulumi.String("batch/v1"),
		Kind:       pulumi.String("Job"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ingress-nginx-admission-create"),
			// Annotations: pulumi.StringMap{
			// "helm.sh/hook":               pulumi.String("pre-install,pre-upgrade"),
			// "helm.sh/hook-delete-policy": pulumi.String("before-hook-creation,hook-succeeded"),
			// },
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		Spec: &batchv1.JobSpecArgs{
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name: pulumi.String("ingress-nginx-admission-create"),
					Labels: pulumi.StringMap{
						// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
						"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
						"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
						"app.kubernetes.io/version":  pulumi.String("0.47.0"),
						// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
						"app.kubernetes.io/component": pulumi.String("admission-webhook"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String("create"),
							Image:           pulumi.String("docker.io/jettech/kube-webhook-certgen:v1.5.1"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Args: pulumi.StringArray{
								pulumi.String("create"),
								pulumi.String(fmt.Sprintf("%v%v%v", "--host=ingress-nginx-controller-admission,ingress-nginx-controller-admission.", "$", "(POD_NAMESPACE).svc")),
								pulumi.String(fmt.Sprintf("%v%v%v", "--namespace=", "$", "(POD_NAMESPACE)")),
								pulumi.String("--secret-name=ingress-nginx-admission"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
							},
						},
					},
					RestartPolicy:      pulumi.String("OnFailure"),
					ServiceAccountName: pulumi.String("ingress-nginx-admission"),
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: pulumi.Bool(true),
						RunAsUser:    pulumi.Int(2000),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = batchv1.NewJob(ctx, "ingress_nginxIngress_nginx_admission_patchJob", &batchv1.JobArgs{
		ApiVersion: pulumi.String("batch/v1"),
		Kind:       pulumi.String("Job"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ingress-nginx-admission-patch"),
			// Annotations: pulumi.StringMap{
			// "helm.sh/hook":               pulumi.String("post-install,post-upgrade"),
			// "helm.sh/hook-delete-policy": pulumi.String("before-hook-creation,hook-succeeded"),
			// },
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		Spec: &batchv1.JobSpecArgs{
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name: pulumi.String("ingress-nginx-admission-patch"),
					Labels: pulumi.StringMap{
						// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
						"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
						"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
						"app.kubernetes.io/version":  pulumi.String("0.47.0"),
						// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
						"app.kubernetes.io/component": pulumi.String("admission-webhook"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String("patch"),
							Image:           pulumi.String("docker.io/jettech/kube-webhook-certgen:v1.5.1"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Args: pulumi.StringArray{
								pulumi.String("patch"),
								pulumi.String("--webhook-name=ingress-nginx-admission"),
								pulumi.String(fmt.Sprintf("%v%v%v", "--namespace=", "$", "(POD_NAMESPACE)")),
								pulumi.String("--patch-mutating=false"),
								pulumi.String("--secret-name=ingress-nginx-admission"),
								pulumi.String("--patch-failure-policy=Fail"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
							},
						},
					},
					RestartPolicy:      pulumi.String("OnFailure"),
					ServiceAccountName: pulumi.String("ingress-nginx-admission"),
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: pulumi.Bool(true),
						RunAsUser:    pulumi.Int(2000),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRole(ctx, "ingress_nginxIngress_nginx_admissionRole", &rbacv1.RoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ingress-nginx-admission"),
			// Annotations: pulumi.StringMap{
			// "helm.sh/hook":               pulumi.String("pre-install,pre-upgrade,post-install,post-upgrade"),
			// "helm.sh/hook-delete-policy": pulumi.String("before-hook-creation,hook-succeeded"),
			// },
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
			Namespace: pulumi.String("ingress-nginx"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "ingress_nginxIngress_nginx_admissionRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("ingress-nginx-admission"),
			// Annotations: pulumi.StringMap{
			// "helm.sh/hook":               pulumi.String("pre-install,pre-upgrade,post-install,post-upgrade"),
			// "helm.sh/hook-delete-policy": pulumi.String("before-hook-creation,hook-succeeded"),
			// },
			Labels: pulumi.StringMap{
				// "helm.sh/chart":                pulumi.String("ingress-nginx-3.33.0"),
				"app.kubernetes.io/name":     pulumi.String("ingress-nginx"),
				"app.kubernetes.io/instance": pulumi.String("ingress-nginx"),
				"app.kubernetes.io/version":  pulumi.String("0.47.0"),
				// "app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/component": pulumi.String("admission-webhook"),
			},
			Namespace: pulumi.String(namespaceNginxIngress.Name),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("Role"),
			Name:     pulumi.String("ingress-nginx-admission"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("ingress-nginx-admission"),
				Namespace: pulumi.String("ingress-nginx"),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
