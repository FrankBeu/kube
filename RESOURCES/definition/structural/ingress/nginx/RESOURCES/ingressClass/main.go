package main

import (
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := networkingv1.NewIngressClass(ctx, "nginxIngressClass", &networkingv1.IngressClassArgs{
			ApiVersion: pulumi.String("networking.k8s.io/v1"),
			Kind:       pulumi.String("IngressClass"),
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"app.kubernetes.io/name":      pulumi.String("ingress-nginx"),
					"app.kubernetes.io/version":   pulumi.String("0.47.0"),
					"app.kubernetes.io/component": pulumi.String("controller"),
				},
				Name: pulumi.String("nginx"),
				Annotations: pulumi.StringMap{
					"ingressclass.kubernetes.io/is-default-class": pulumi.String("true"),
				},
			},
			Spec: &networkingv1.IngressClassSpecArgs{
				Controller: pulumi.String("k8s.io/ingress-nginx"),
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
