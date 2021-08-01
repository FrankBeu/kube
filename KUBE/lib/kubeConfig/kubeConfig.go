// Package config contains stackConfiguration-related helpers
package kubeConfig

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// KubeConfig makes it possible to control resource creation
// by commenting or uncommenting creators
type (
	KubeConfig map[string]func(*pulumi.Context) error
)

func Env(ctx *pulumi.Context) (env string) {
	conf := config.New(ctx, "")
	env = conf.Require("env")
	return
}

type grafanaSecret struct {
	AdminPassword string
}

func GrafanaAdminPassword(ctx *pulumi.Context) (grafanaAdminPw string) {
	var gs grafanaSecret
	conf := config.New(ctx, "")
	conf.RequireSecretObject("grafanaSecret", &gs)

	grafanaAdminPw = gs.AdminPassword

	return
}

type domainSecret struct {
	AdminEmail string
	Domain     string
}

func AdminEmail(ctx *pulumi.Context) (adminEmail string) {
	var ds domainSecret
	conf := config.New(ctx, "")
	conf.RequireSecretObject("domainSecret", &ds)

	adminEmail = ds.AdminEmail

	return
}

// DomainNameSuffix returns the string ".DOMAIN-DEFINED-IN-PULUMI-CONFIG" - appendable to subDomainNames
func DomainNameSuffix(ctx *pulumi.Context) (domainNameSuffix string) {
	var ds domainSecret
	conf := config.New(ctx, "")
	conf.RequireSecretObject("domainSecret", &ds)
	domainNameSuffix = "." + ds.Domain
	return
}

// func DomainClusterIssuer(ctx *pulumi.Context) (clusterIssuer int) {

// 	var ds domainSecret
// 	conf := config.New(ctx, "")
// 	conf.RequireSecretObject("domainSecret", &ds)
// 	domainClusterIssuer =   ds.ClusterIssuer


// 	return
// }
