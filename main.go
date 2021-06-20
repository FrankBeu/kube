package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/definition/structural/ingress/nginx"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		////////////////////////
		// STRUCTURAL
		//
		err := nginx.CreateNginxIngressController(ctx)
		return err


		////////////////////////
		// ALT
		//
		// return nil
	})
}
