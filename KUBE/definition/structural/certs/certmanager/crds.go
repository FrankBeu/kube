package certmanager

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Crd struct {
	Name     string
	Location string
}

func registerCRDs(ctx *p.Context, crds []Crd) error {

	for _, crd := range crds {

		_, err := yaml.NewConfigFile(ctx, crd.Name,
			&yaml.ConfigFileArgs{
				File: crd.Location,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
