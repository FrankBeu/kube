// Package petstore is a testApp for the glooEdge
package petstore

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/lib"
)

func CreateGlooPetstore(ctx *p.Context) error {

	namespacePetstore := &lib.Namespace{
		Name:          "testing-petstore",
		Tier:          lib.NamespaceTierTesting,
		GlooDiscovery: true,
	}

	err := lib.CreateNamespace(ctx, namespacePetstore)

	if err != nil {
		return err
	}

	_, err = appsv1.NewDeployment(ctx, "defaultPetstoreDeployment", &appsv1.DeploymentArgs{
		ApiVersion: p.String("apps/v1"),
		Kind:       p.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app": p.String("petstore"),
			},
			Name:      p.String("petstore"),
			Namespace: p.String(namespacePetstore.Name),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: p.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: p.StringMap{
					"app": p.String("petstore"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: p.StringMap{
						"app": p.String("petstore"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Image: p.String("soloio/petstore-example:latest"),
							Name:  p.String("petstore"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: p.Int(8080),
									Name:          p.String("http"),
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
		ApiVersion: p.String("v1"),
		Kind:       p.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"service": p.String("petstore"),
			},
			Name:      p.String("petstore"),
			Namespace: p.String(namespacePetstore.Name),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:     p.Int(8080),
					Protocol: p.String("TCP"),
				},
			},
			Selector: p.StringMap{
				"app": p.String("petstore"),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
