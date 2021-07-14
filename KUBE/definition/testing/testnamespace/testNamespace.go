// Package testnamespace defines the namespace used for tests, protos ...
// TODO: Create as singleton; delegate creation to object requiring the namespace
package testnamespace

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/namespace"
)

var (
	name          = "test"
	namespaceTest = &namespace.Namespace{
		Name: name,
		Tier: namespace.NamespaceTierTesting,
		// GlooDiscovery: true,
	}
)

func CreateTestNamespace(ctx *pulumi.Context) error {
	_, err := namespace.CreateNamespace(ctx, namespaceTest)
	if err != nil {
		return err
	}

	return nil
}
