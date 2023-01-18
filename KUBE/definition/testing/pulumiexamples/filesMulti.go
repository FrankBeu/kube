package pulumiexamples

import (
	"path/filepath"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateFromFilesMulti(ctx *pulumi.Context) error {
	guestbook, err := yaml.NewConfigGroup(ctx, "guestbook", &yaml.ConfigGroupArgs{
		Files: []string{filepath.Join("definition/testing/pulumiexamples/fileSingle", "*.yaml")}, //nolint:gocritic
	})
	if err != nil {
		return err
	}

	// Export the private cluster IP address of the frontend.
	// can be retrieved with `pulumi stack output privateIP`
	frontend := guestbook.GetResource("v1/Service", "frontend", "").(*corev1.Service)
	ctx.Export("fileMulti.privateIP", frontend.Spec.ClusterIP()) //// retrievable with `p stack output fileMulti.privateIP`

	return nil
}
