package testutil

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Mocks int

func (Mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (Mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

// TestConfig ist the default stackConfiguration used in tests
var TestConfig = map[string]string{
	"project:domain": "test.domain",
}

// WithMocksAndConfig enables running unit tests with a test-stackConfiguration (Pulumi.kube-….yaml)
func WithMocksAndConfig(project, stack string, config map[string]string, mocks pulumi.MockResourceMonitor) pulumi.RunOption {
	return func(info *pulumi.RunInfo) {
		info.Project, info.Stack, info.Mocks, info.Config = project, stack, mocks, config
	}
}
