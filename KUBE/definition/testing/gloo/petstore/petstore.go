// Package petstore is a testApp for the glooEdge
package petstore

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"thesym.site/kube/lib/namespace"
)

var (
	name              = "petstore"
	namespacePetstore = &namespace.Namespace{
		Name:          name,
		Tier:          namespace.NamespaceTierTesting,
		GlooDiscovery: true,
	}
)

func CreateGlooPetstore(ctx *pulumi.Context) error {

	conf := config.New(ctx, "")
	domainName := name + "." + conf.Require("domain")

	_, err := namespace.CreateNamespace(ctx, namespacePetstore)
	if err != nil {
		return err
	}

	err = addIngress(ctx, domainName)

	if err != nil {
		return err
	}

	err = execGeneratedCode(ctx)

	if err != nil {
		return err
	}

	return nil
}

// TODO: extract
func addIngress(ctx *pulumi.Context, domainName string) error {

	_, err := networkingv1.NewIngress(ctx, namespacePetstore.Name, &networkingv1.IngressArgs{
		ApiVersion: pulumi.String("networking.k8s.io/v1"),
		Kind:       pulumi.String("Ingress"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels:      pulumi.StringMap{"app": pulumi.String("petstore")},
			Name:        pulumi.String("petstore"),
			Namespace:   pulumi.String(namespacePetstore.Name),
			Annotations: pulumi.StringMap{
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
			Rules: networkingv1.IngressRuleArray{
				&networkingv1.IngressRuleArgs{
					Host: pulumi.String(domainName),
					Http: &networkingv1.HTTPIngressRuleValueArgs{
						Paths: networkingv1.HTTPIngressPathArray{
							&networkingv1.HTTPIngressPathArgs{
								PathType: pulumi.String("Prefix"),
								Path:     pulumi.String("/"),
								Backend: &networkingv1.IngressBackendArgs{
									Service: &networkingv1.IngressServiceBackendArgs{
										Name: pulumi.String("petstore"),
										Port: &networkingv1.ServiceBackendPortArgs{
											Number: pulumi.Int(8080),
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

func execGeneratedCode(ctx *pulumi.Context) error {
	// CHANGES: namespaceAdded

	_, err := appsv1.NewDeployment(ctx, "defaultPetstoreDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app": pulumi.String("petstore"),
			},
			Name:      pulumi.String("petstore"),
			Namespace: pulumi.String(namespacePetstore.Name),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String("petstore"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String("petstore"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Image: pulumi.String("soloio/petstore-example:latest"),
							Name:  pulumi.String("petstore"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8080),
									Name:          pulumi.String("http"),
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

	_, err = corev1.NewService(ctx, "defaultPetstoreService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"service": pulumi.String("petstore"),
			},
			Name:      pulumi.String("petstore"),
			Namespace: pulumi.String(namespacePetstore.Name),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:     pulumi.Int(8080),
					Protocol: pulumi.String("TCP"),
				},
			},
			Selector: pulumi.StringMap{
				"app": pulumi.String("petstore"),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
