// Package kubeconfig contains stackConfiguration-related helpers
package kubeconfig

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"thesym.site/kube/lib/types"
)

// KubeConfig makes it possible to control resource creation
// by commenting or uncommenting creators
type (
	KubeConfig map[string]func(*pulumi.Context) error
)

type grafanaSecret struct {
	AdminPassword string
}

type domainSecret struct {
	AdminEmail string
	Domain     string
}

type domain struct {
	ClusterIssuer string
}

func Env(ctx *pulumi.Context) (env string) {
	conf := config.New(ctx, "")
	env = conf.Require("env")
	return
}

func GrafanaAdminPassword(ctx *pulumi.Context) (grafanaAdminPw string) {
	var gs grafanaSecret
	conf := config.New(ctx, "")
	conf.RequireSecretObject("grafanaSecret", &gs)

	grafanaAdminPw = gs.AdminPassword

	return
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

func DomainClusterIssuer(ctx *pulumi.Context) (clusterIssuerType types.ClusterIssuerType) {
	var d domain
	conf := config.New(ctx, "")
	conf.RequireObject("domain", &d)
	domainClusterIssuer := d.ClusterIssuer
	clusterIssuerType, err := convertClusterIssuerTypeString2Type(domainClusterIssuer)
	if err != nil {
		panic(err)
	}

	return
}

func convertClusterIssuerTypeString2Type(clusterIssuerTypeString string) (types.ClusterIssuerType, error) {
	for key, cit := range types.AllClusterIssuerTypes {
		if clusterIssuerTypeString == cit.String() {
			return types.ClusterIssuerType(key), nil
		}
	}
	return 0, errors.New("could not find clusterIssuerType - check your pulumiConfigFiles")
}
