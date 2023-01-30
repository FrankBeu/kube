package whoami

import (
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/lib/ingressroute"
	"thesym.site/kube/lib/types"
)

var (
	subDomainNameIngressRoute = "ingressroute." + name
	whoamiRouteConf           = types.IngressRouteConfig{
		Service: types.Service{
			Name:          name,
			NamespaceName: testnamespace.NamespaceTest.Name,
			Port:          80,
			Kind:          types.ServiceKindK8S,
		},
		Name:              strings.ReplaceAll(subDomainNameIngressRoute, ".", "-"),
		HostSubdomainName: subDomainNameIngressRoute,
		NamespaceName:     namespace,
	}
)

func CreateIngressRoute(ctx *pulumi.Context) error {
	_, err := ingressroute.CreateIngressRoute(ctx, &whoamiRouteConf)
	if err != nil {
		return err
	}

	return nil
}
