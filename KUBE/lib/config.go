package lib

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// KubeConfig makes it possible to control which resources will be installed
// by commenting or uncommenting creators
type (
	KubeConfig map[string]func(*pulumi.Context) error
)
