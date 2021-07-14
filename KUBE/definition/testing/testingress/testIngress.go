// Package testingress use used for testing tls-ingresses
// TODO: cleanup and extractions required
package testingress

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"thesym.site/kube/lib/certificate"
)

// CreateTestIngress creates an application with a tls-ingress
func CreateTestIngress(ctx *pulumi.Context) error {
	namespace := "test"
	name := "testingress"

	err := createDeploymentAndService(ctx, name, namespace)
	if err != nil {
		return err
	}

	err = addIngress(ctx, name, namespace)
	if err != nil {
		return err
	}

	return nil
}

// TODO: extract (gloopetstore, too)
func addIngress(ctx *pulumi.Context, name, namespace string) error {
	conf := config.New(ctx, "")
	domainNameSuffix := "." + conf.Require("domain")

	_, err := networkingv1.NewIngress(ctx, name, &networkingv1.IngressArgs{
		ApiVersion: pulumi.String("networking.k8s.io/v1"),
		Kind:       pulumi.String("Ingress"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels:    pulumi.StringMap{"app": pulumi.String(name)},
			Name:      pulumi.String(name),
			Namespace: pulumi.String(namespace),
			Annotations: pulumi.StringMap{
				//// TLS
				"cert-manager.io/cluster-issuer": pulumi.String(certificate.ClusterIssuerTypeCaLocal.String()),
				//// TLS-END
				// "kubernetes.io/ingress.class": pulumi.String("nginx"),
				// "nginx.ingress.kubernetes.io/force-ssl-redirect": pulumi.String("true"),
				// "nginx.ingress.kubernetes.io/ssl-redirect": pulumi.String("true"),
				// "nginx.ingress.kubernetes.io/ssl-redirect": pulumi.String("false"),
				// "nginx.ingress.kubernetes.io/configuration-snippet": pulumi.String(
				// 	`if ($http_x_forwarded_proto = 'http') {
				//            return 301 https://$host$request_uri;
				//          }`),
			},
		},
		Spec: &networkingv1.IngressSpecArgs{
			//// TODO: create const
			IngressClassName: pulumi.String("nginx"),
			//// TLS

			Tls: networkingv1.IngressTLSArray{
				&networkingv1.IngressTLSArgs{
					Hosts: pulumi.StringArray{
						pulumi.String(name + domainNameSuffix),
					},
					SecretName: pulumi.String(name + "-tls"),
				},
			},
			//// TLS-END
			Rules: networkingv1.IngressRuleArray{
				&networkingv1.IngressRuleArgs{
					Host: pulumi.String(name + domainNameSuffix),
					Http: &networkingv1.HTTPIngressRuleValueArgs{
						Paths: networkingv1.HTTPIngressPathArray{
							&networkingv1.HTTPIngressPathArgs{
								PathType: pulumi.String("Prefix"),
								Path:     pulumi.String("/"),
								Backend: &networkingv1.IngressBackendArgs{
									Service: &networkingv1.IngressServiceBackendArgs{
										Name: pulumi.String(name),
										Port: &networkingv1.ServiceBackendPortArgs{
											Number: pulumi.Int(80),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// Tls:              nil,
	})
	if err != nil {
		return err
	}
	return nil
}

func createDeploymentAndService(ctx *pulumi.Context, name string, namespace string) error {
	appName := name
	appLabels := pulumi.StringMap{
		"app": pulumi.String(appName),
	}
	_, err := appsv1.NewDeployment(ctx, name, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels:    appLabels,
			Namespace: pulumi.String(namespace),
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Replicas: pulumi.Int(1),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String(appName),
							Image: pulumi.String("nginx:1.15-alpine"),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = corev1.NewService(ctx, appName, &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Labels:    appLabels,
			Namespace: pulumi.String(namespace),
			Name:      pulumi.String(appName),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type: pulumi.String("ClusterIP"),
			Ports: &corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:       pulumi.Int(80),
					TargetPort: pulumi.Int(80),
					Protocol:   pulumi.String("TCP"),
				},
			},
			Selector: appLabels,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
