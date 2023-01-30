// Package ingressroute is responsible for the creation of ingressroute-related resources.
package ingressroute

import (
	"regexp"
	"strconv"
	"strings"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	certv1 "thesym.site/kube/crds/kubernetes/certmanager/v1"
	traefikCRD "thesym.site/kube/crds/kubernetes/traefik/v1alpha1"

	// "thesym.site/kube/definition/structural/ingress/traefik"
	"thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/kubeconfig"
	"thesym.site/kube/lib/types"
)

// CreateIngressRoute creates an ingressRoute and the required certificate;
// ensuring that the domain used as ir.match and the crt.commonName are the same and the
// ir.Tls.SecretName and cert.secretName are the same, too.
func CreateIngressRoute(ctx *pulumi.Context, ingrtcnf *types.IngressRouteConfig) (*traefikCRD.IngressRoute, error) {
	domainNameSuffix := kubeconfig.DomainNameSuffix(ctx)

	_, ingressRoute, err := createIngressRoute(ctx, ingrtcnf, domainNameSuffix)
	if err != nil {
		return nil, err
	}

	//// export the accessible url exposed by the ingressroute
	//// currently only one route is used per ingressroute
	url := ingressRoute.Spec.Routes().Index(pulumi.Int(0)).Match().ApplyT(extractURLFromMatch)
	ctx.Export("URL."+ingrtcnf.Name, url)

	return ingressRoute, nil
}

//nolint:lll
func createIngressRoute(ctx *pulumi.Context, ingrtcnf *types.IngressRouteConfig, domainNameSuffix string) (*certv1.Certificate, *traefikCRD.IngressRoute, error) {
	//// match overrides ir.match and c.commonName
	match := ""
	commonNameSegment := ""
	if ingrtcnf.Match != "" {
		match = ingrtcnf.Match
		extractedName := strings.ReplaceAll(string(extractURLFromMatch(match)), kubeconfig.DomainNameSuffix(ctx), "")
		commonNameSegment = extractedName
	} else {
		match = "Host(`" + ingrtcnf.HostSubdomainName + domainNameSuffix + "`)"
		commonNameSegment = ingrtcnf.HostSubdomainName
	}

	ingressRouteCert := types.Cert{
		ClusterIssuerType:        kubeconfig.DomainClusterIssuer(ctx),
		Namespace:                ingrtcnf.NamespaceName,
		Name:                     ingrtcnf.Name,
		CommonNameSegment:        commonNameSegment,
		Duration:                 strconv.Itoa(types.IngressRouteTLSDefaultDurationInDays*24) + "h", //nolint:gomnd
		AdditionalSubdomainNames: []string{},
	}

	cert, err := certificate.CreateCert(ctx, &ingressRouteCert)
	if err != nil {
		return nil, nil, err
	}

	//// middleware references only
	middlewares := traefikCRD.IngressRouteSpecRoutesMiddlewaresArray{}
	if len(ingrtcnf.Middlewares) > 0 {
		for _, m := range ingrtcnf.Middlewares {
			middleware := traefikCRD.IngressRouteSpecRoutesMiddlewaresArgs{
				Name: pulumi.String(m.Name.String()),
				// Namespace: pulumi.String(traefik.NamespaceTraefikIngressController.Name),
				Namespace: pulumi.String("ingress-traefik"), ////TODO cannot import traefik
			}
			middlewares = append(middlewares, middleware)
		}
	}

	////TODO test
	////a traefikService does not need a port and a namespace
	services := traefikCRD.IngressRouteSpecRoutesServicesArray{}
	service := traefikCRD.IngressRouteSpecRoutesServicesArgs{
		Kind: pulumi.String(ingrtcnf.Service.Kind.String()),
		Name: pulumi.String(ingrtcnf.Service.Name),
	}

	if ingrtcnf.Service.Kind == types.ServiceKindK8S {
		service.Namespace = pulumi.String(ingrtcnf.Service.NamespaceName)
		service.Port = pulumi.Int(ingrtcnf.Service.Port)
	}
	services = append(services, service)

	ingRouteSec, err := traefikCRD.NewIngressRoute(
		ctx,
		ingrtcnf.Name,
		&traefikCRD.IngressRouteArgs{
			ApiVersion: pulumi.String("traefik.containo.us/v1alpha1"),
			Kind:       pulumi.String("ingressRoute"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(ingrtcnf.Name),
				Namespace: pulumi.String(ingrtcnf.NamespaceName),
			},
			Spec: &traefikCRD.IngressRouteSpecArgs{
				EntryPoints: pulumi.StringArray{
					pulumi.String("websecure"),
				},
				Routes: traefikCRD.IngressRouteSpecRoutesArray{
					traefikCRD.IngressRouteSpecRoutesArgs{
						Kind:  pulumi.String("Rule"),
						Match: pulumi.String(match),

						Middlewares: middlewares,
						Services:    services,
					},
				},
				Tls: &traefikCRD.IngressRouteSpecTlsArgs{
					////has to fit the cert.spec.SecretName and must not contain "."s
					SecretName: pulumi.String(ingrtcnf.Name + types.IngressRouteTLSSecretSuffix),
				},
			},
		})

	if err != nil {
		return nil, nil, err
	}

	return cert, ingRouteSec, nil
}

func extractURLFromMatch(match string) pulumi.String {
	re := regexp.MustCompile(`^Host\(.(.*?).\)`)
	url := re.FindStringSubmatch(match)[1]
	return pulumi.String(url)
}
