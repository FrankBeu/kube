package lib

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Crd defines the necessary fields to handle a CustomResourceDefinition provided by `crd2pulumi`
type Crd struct {
	Name string
	// The Location must be relative to the go-module-root.
	// Otherwise only the kubernetes:yaml:ConfigFile will be created
	// kubernetes:...:CustomResourceDefinition are needed, too
	// multiple CRDs can be defined in one file
	Location string
}

func RegisterCRDs(ctx *p.Context, crds []Crd) error {
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
