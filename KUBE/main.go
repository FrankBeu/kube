package main

import (
	"strconv"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiConfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	empty "thesym.site/kube/env/auxiliary/emptycluster"
	full "thesym.site/kube/env/auxiliary/fullcluster"
	dev "thesym.site/kube/env/development"
	prod "thesym.site/kube/env/production"
	stage "thesym.site/kube/env/staging"
	"thesym.site/kube/lib/config"
	"thesym.site/kube/lib/debug"
)

// for convenience - preUse all environments
var (
	_ = empty.Kube
	_ = full.Kube
	_ = dev.Kube
	_ = prod.Kube
	_ = stage.Kube
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := pulumiConfig.New(ctx, "")

		//// DebugMode cf. [[file:../README.org::*debugging with delve]]
		debugMode, _ := strconv.ParseBool(conf.Get("debugMode"))
		timeOutDuration := 120 * time.Second
		debugReady := false

		err := debug.RunInDebugMode(debugMode, timeOutDuration, debugReady)
		if err != nil {
			return err
		}

		var kube config.KubeConfig

		//// Load the configuration for the current stack
		switch env := conf.Require("env"); env {
		case "dev":
			// kube = empty.Kube
			// kube = full.Kube
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
