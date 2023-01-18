// Package testingress use used for testing tls-ingresses
// TODO: cleanup and extractions required
package testingress

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/ingress"
	"thesym.site/kube/lib/types"
)

// CreateTestIngress creates an application with a tls-ingress
func CreateTestIngress(ctx *pulumi.Context) error {
	namespace := "test"
	name := "testingress"
	servicePort := 80
	serviceTargetPort := 80

	err := createDeploymentAndService(ctx, name, namespace, servicePort, serviceTargetPort)
	if err != nil {
		return err
	}

	ingConf := types.IngressConfig{
		Annotations: map[string]pulumi.StringInput{
			"traefik.ingress.kubernetes.io/router.entrypoints": pulumi.String("websecure"),
			"pulumi.com/skipAwait":                             pulumi.String("true"),
		},
		ClusterIssuerType: types.ClusterIssuerTypeCALocal,
		Hosts: []types.Host{
			{
				Name:        name,
				ServiceName: name,
				ServicePort: servicePort,
			},
		},
		IngressClassName: types.IngressClassNameTraefik,
		Name:             name,
		NamespaceName:    namespace,
		TLS:              true,
	}
	testIng, err := ingress.CreateIngress(ctx, &ingConf)
	if err != nil {
		return err
	}
	// _ = ingConf
	_ = testIng
	// ctx.Export("testIngressUrl", testIng.Spec.Rules().Index(pulumi.Int(0)).Host()) //// retrievable with `p stack output tesIngressUrl`

	return nil
}

func createDeploymentAndService(ctx *pulumi.Context, name, namespace string, servicePort, serviceTargetPort int) error {
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
					Port:       pulumi.Int(servicePort),
					TargetPort: pulumi.Int(serviceTargetPort),
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
