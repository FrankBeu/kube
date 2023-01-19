// Package pulumiexamples contains all the introductory examples for references
package pulumiexamples

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateFromFileSingle(ctx *pulumi.Context) error {
	guestbook, err := yaml.NewConfigFile(ctx, "guestbook", &yaml.ConfigFileArgs{
		// _, err := yaml.NewConfigFile(ctx, "guestbook", &yaml.ConfigFileArgs{
		//nolint:lll
		// File: "fileSingle/guestbook-all-in-one.yaml",                                //// will only create  └─ kubernetes:yaml:ConfigFile guestbook
		File: "definition/testing/pulumiexamples/fileSingle/guestbook-all-in-one.yaml", //// will create the ConfigFile + the Resources

	})
	if err != nil {
		return err
	}

	// // Export the private cluster IP address of the frontend.
	frontend := guestbook.GetResource("v1/Service", "frontend", "").(*corev1.Service)
	ctx.Export("fileSingle.privateIP:", frontend.Spec.ClusterIP()) //// retrievable with `p stack output fileSingle.privateIP`

	return nil
}
