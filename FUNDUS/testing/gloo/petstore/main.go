package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := appsv1.NewDeployment(ctx, "defaultPetstoreDeployment", &appsv1.DeploymentArgs{
			ApiVersion: pulumi.String("apps/v1"),
			Kind:       pulumi.String("Deployment"),
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"app": pulumi.String("petstore"),
				},
				Name:      pulumi.String("petstore"),
				Namespace: pulumi.String("default"),
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
				Namespace: pulumi.String("default"),
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
	})
}
