package testhelmrelease

import (
	"fmt"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	helm "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateHelmRelease(ctx *pulumi.Context) error {
	namespace := "test"
	// https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/kubernetes-go-helm-release-wordpress/
	// https://github.com/pulumi/examples/tree/master/kubernetes-go-helm-release-wordpress

	wordpress, err := helm.NewRelease(ctx, "wpdev", &helm.ReleaseArgs{
		Version:   pulumi.String("15.2.5"),
		Chart:     pulumi.String("wordpress"),
		Namespace: pulumi.String(namespace),
		Values: pulumi.Map{
			"service": pulumi.StringMap{
				"type": pulumi.String("ClusterIP"),
			},
		},
		RepositoryOpts: &helm.RepositoryOptsArgs{
			Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
		},
	})
	if err != nil {
		return err
	}

	//// Await on the Status field of the wordpress release and use that to lookup the WordPress service.
	result := pulumi.All(wordpress.Status.Namespace(), wordpress.Status.Name()).ApplyT(func(r interface{}) ([]interface{}, error) {
		arr := r.([]interface{})
		namespace := arr[0].(*string)
		name := arr[1].(*string)
		svc, err := corev1.GetService(ctx, "svc", pulumi.ID(fmt.Sprintf("%s/%s-wordpress", *namespace, *name)), nil)
		if err != nil {
			return nil, err
		}
		//// Return the cluster IP and service name
		return []interface{}{svc.Spec.ClusterIP().Elem(), svc.Metadata.Name().Elem()}, nil
	})

	arr := result.(pulumi.ArrayOutput)
	ctx.Export("testHelmRelease.frontendIp", arr.Index(pulumi.Int(0)))
	ctx.Export("testHelmRelease.portForwardCommand", pulumi.Sprintf("kubectl port-forward svc/%s 8080:80", arr.Index(pulumi.Int(1))))
	return nil
}
