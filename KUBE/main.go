package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	dev "thesym.site/kube/env/development"
	prod "thesym.site/kube/env/production"
	stage "thesym.site/kube/env/staging"
	"thesym.site/kube/lib"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")

		var kube lib.KubeConfig

		switch env := conf.Require("env"); env {
		case "dev":
			kube = dev.Kube
		case "stage":
			kube = stage.Kube
		case "prod":
			kube = prod.Kube
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
