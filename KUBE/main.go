package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/definition/testing/gloo/petstore"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		////////////////////////
		// STRUCTURAL
		//
		// err := nginx.CreateNginxIngressController(ctx)
		// return err

		////////////////////////
		// TESTING
		//
		err := petstore.CreateGlooPetstore(ctx)
		return err

		////////////////////////
		// ALT
		//
		// return nil
	})
}
