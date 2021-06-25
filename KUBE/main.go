package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/definition/structural/certs/certmanager"
	"thesym.site/k8s/definition/structural/ingress/nginx"
	"thesym.site/k8s/definition/testing/gloo/petstore"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		//// config
		kube := map[string]func(*pulumi.Context) error{
			//////////////////////// //////////////////////// ////////////////////////
			//// STRUCTURAL
			////
			//////////////////////// ////////////////////////
			//// Ingress
			////
			"nginxIngress": nginx.CreateNginxIngressController,

			//// installation working crd-usage not
			// "emmissary": emmissary.CreateEmmissary,

			// NOTREADY gloo

			// "tykGateway": tyk.CreateTykGateway,

			//////////////////////// ////////////////////////
			//// certificates
			////
			"certmanager": certmanager.CreateCertmanager,

			//////////////////////// //////////////////////// ////////////////////////
			//// TESTING
			////
			"glooPetstore": petstore.CreateGlooPetstore,
		}

		for _, creator := range kube {
			err := creator(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
