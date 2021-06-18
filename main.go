package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/k8s/structural/ingress/nginx"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := nginx.CreateNginxIngressController(ctx)
		return err

	})
}
