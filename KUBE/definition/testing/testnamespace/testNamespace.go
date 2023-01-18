// Package testnamespace defines the namespace used for tests, protos ...
// TODO: Create as singleton; delegate creation to object requiring the namespace
package testnamespace

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/namespace"
)

var (
	name          = "test"
	NamespaceTest = &namespace.Namespace{
		Name: name,
		Tier: namespace.NamespaceTierTesting,
	}
)

func CreateTestNamespace(ctx *pulumi.Context) error {
	_, err := namespace.CreateNamespace(ctx, NamespaceTest)
	if err != nil {
		return err
	}

	return nil
}
