// Package config contains stackConfiguration-related helpers
package config

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// KubeConfig makes it possible to control resource creation
// by commenting or uncommenting creators
type (
	KubeConfig map[string]func(*pulumi.Context) error
)

// DomainNameSuffix returns the string ".DOMAIN-DEFINED-IN-PULUMI-CONFIG" - appendable to subDomainNames
func DomainNameSuffix(ctx *pulumi.Context) (domainNameSuffix string) {
	conf := config.New(ctx, "")
	domainNameSuffix = "." + conf.Require("domain")
	return
}

func Env(ctx *pulumi.Context) (env string) {
	conf := config.New(ctx, "")
	env = conf.Require("env")
	return
}

type grafana struct {
	AdminPassword string
}

func GrafanaAdminPassword(ctx *pulumi.Context) (grafanaAdminPw string) {
	var g grafana
	conf := config.New(ctx, "")
	conf.RequireSecretObject("grafana", &g)

	grafanaAdminPw = g.AdminPassword

	return
}
