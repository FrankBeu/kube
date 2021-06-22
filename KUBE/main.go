package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/definition/structural/certs/certmanager"
	"thesym.site/k8s/definition/structural/ingress/tyk"
	"thesym.site/k8s/definition/testing/gloo/petstore"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		//////////////////////// //////////////////////// ////////////////////////
		// STRUCTURAL
		//

		//////////////////////// ////////////////////////
		// Ingress

		////////////////////////
		// nginx
		// err := nginx.CreateNginxIngressController(ctx)
		// if err != nil {
		// return err
		// }

		////////////////////////
		// ambassador

		////////////////////////
		// gloo

		////////////////////////
		// tykGateway
		err := tyk.CreateTykGateway(ctx)
		if err != nil {
			return err
		}

		//////////////////////// ////////////////////////
		// certificates
		////////////////////////
		// certManager
		err = certmanager.CreateCertmanager(ctx)
		if err != nil {
			return err
		}

		//////////////////////// //////////////////////// ////////////////////////
		// TESTING
		//

		////////////////////////
		// glooPetstore
		err = petstore.CreateGlooPetstore(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
