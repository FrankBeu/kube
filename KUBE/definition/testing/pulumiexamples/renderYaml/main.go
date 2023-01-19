package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		renderProvider, err := kubernetes.NewProvider(ctx, "k8s-yaml-renderer", &kubernetes.ProviderArgs{
			RenderYamlToDirectory: pulumi.String("yaml"),
		})
		if err != nil {
			return err
		}

		labels := pulumi.StringMap{"app": pulumi.String("nginx")}

		_, err = appsv1.NewDeployment(ctx, "nginx-dep", &appsv1.DeploymentArgs{
			Spec: &appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{MatchLabels: labels},
				Replicas: pulumi.Int(1),
				Template: corev1.PodTemplateSpecArgs{
					Metadata: metav1.ObjectMetaArgs{Labels: labels},
					Spec: &corev1.PodSpecArgs{
						Containers: &corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("nginx"),
								Image: pulumi.String("nginx"),
							},
						},
					},
				},
			},
		}, pulumi.Provider(renderProvider))
		if err != nil {
			return err
		}
		_, err = corev1.NewService(ctx, "nginx-svc", &corev1.ServiceArgs{
			Spec: &corev1.ServiceSpecArgs{
				Type:     pulumi.String("LoadBalancer"),
				Selector: labels,
				Ports:    corev1.ServicePortArray{corev1.ServicePortArgs{Port: pulumi.Int(80)}}, //nolint:gomnd
			},
		}, pulumi.Provider(renderProvider))
		if err != nil {
			return err
		}

		return nil
	})
}
