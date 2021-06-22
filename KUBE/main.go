package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/definition/structural/certs/certmanager"
	"thesym.site/k8s/definition/testing/gloo/petstore"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		//////////////////////// //////////////////////// ////////////////////////
		// STRUCTURAL
		//

		//////////////////////// ////////////////////////
		// IngressController
		////////////////////////
		// nginx
		// err := nginx.CreateNginxIngressController(ctx)
		// if err != nil {
		// return err
		// }

		////////////////////////
		// gloo

		//////////////////////// ////////////////////////
		// certificates
		////////////////////////
		// certManager
		err := certmanager.CreateCertmanager(ctx)
		if err != nil {
			return err
		}

		//////////////////////// //////////////////////// ////////////////////////
		// TESTING
		//
		err = petstore.CreateGlooPetstore(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
