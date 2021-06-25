package pulumiexamples

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateFromFileSingleWithConfigurationTransformation(ctx *pulumi.Context) error {
	guestbook, err := yaml.NewConfigFile(ctx, "guestbook", &yaml.ConfigFileArgs{
		// _, err := yaml.NewConfigFile(ctx, "guestbook", &yaml.ConfigFileArgs{
		File: "definition/testing/pulumiexamples/fileSingle/guestbook-all-in-one.yaml",
		Transformations: []yaml.Transformation{func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
			if kind, ok := state["kind"]; ok && kind == "Service" && state["metadata"].(map[string]interface{})["name"] == "frontend" {
				state["spec"].(map[string]interface{})["type"] = "LoadBalancer"
			}
		}},
	})
	if err != nil {
		return err
	}

	// Export the private cluster IP address of the frontend.
	frontend := guestbook.GetResource("v1/Service", "frontend", "").(*corev1.Service)
	ctx.Export("privateIP", frontend.Spec.ClusterIP())
	ctx.Export("publicIP", frontend.Status.LoadBalancer().Ingress().Index(pulumi.Int(0)).Ip())

	return nil
}
