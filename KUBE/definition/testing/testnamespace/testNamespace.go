// Package testnamespace defines the namespace used for tests, protos ...
// TODO: Create as singleton; delegate creation to object requiring the namespace
package testnamespace

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib"
)

var (
	name          = "test"
	namespaceTest = &lib.Namespace{
		Name: name,
		Tier: lib.NamespaceTierTesting,
		// GlooDiscovery: true,
	}
)

func CreateTestNamespace(ctx *pulumi.Context) error {
	err := lib.CreateNamespaces(ctx, namespaceTest)
	if err != nil {
		return err
	}

	return nil
}
