// Package whoami defines a whoami-app
package whoami

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/definition/testing/testnamespace"
)

var (
	name        = "whoami"
	namespace   = testnamespace.NamespaceTest.Name
	defaultPort = 80
)

func CreateWhoAmI(ctx *pulumi.Context) error {
	_, err := corev1.NewService(ctx, "whoamiService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("whoami"),
			Namespace: pulumi.String(namespace),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("http"),
					TargetPort: pulumi.Any(defaultPort),
					Port:       pulumi.Int(defaultPort),
					Protocol:   pulumi.String("TCP"),
				},
			},
			Selector: pulumi.StringMap{
				"app": pulumi.String("whoami"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "whoamiDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("whoami-deployment"),
			Namespace: pulumi.String(testnamespace.NamespaceTest.Name),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String("whoami"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String("whoami"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("whoami-container"),
							Image: pulumi.String("containous/whoami"),
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
