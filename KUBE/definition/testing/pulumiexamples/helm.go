package pulumiexamples

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateFromHelmChart(ctx *pulumi.Context) error {

	// Deploy the latest version of the stable/wordpress chart.
	wordpress, err := helm.NewChart(ctx, "wpdev", helm.ChartArgs{
		Repo:    pulumi.String("stable"),
		Chart:   pulumi.String("wordpress"),
		Version: pulumi.String("9.0.3"),
	})
	if err != nil {
		return err
	}

	// Export the public IP for WordPress.
	// frontendIP := wordpress.GetResource("v1/Service", "wpdev-wordpress", "").Apply(func(r interface{}) (interface{}, error) {
	frontendIP := wordpress.GetResource("v1/Service", "wpdev-wordpress", "").ApplyT(func(r interface{}) (interface{}, error) {
		svc := r.(*corev1.Service)
		return svc.Status.LoadBalancer().Ingress().Index(pulumi.Int(0)).Ip(), nil
	})
	ctx.Export("helmChart.frontendIp", frontendIP)

	return nil
}
