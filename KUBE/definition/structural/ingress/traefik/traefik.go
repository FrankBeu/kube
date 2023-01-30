// Package traefik provides the default traefikIngressController-installation
package traefik

import (
	helm "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	traefikCRD "thesym.site/kube/crds/kubernetes/traefik/v1alpha1"
	"thesym.site/kube/lib/ingressroute"
	"thesym.site/kube/lib/kubeconfig"
	"thesym.site/kube/lib/namespace"
	"thesym.site/kube/lib/types"
)

var (
	ingressControllerPortHTTP         = 30080
	ingressControllerPortHTTPS        = 30443
	NamespaceTraefikIngressController = &namespace.Namespace{
		Name: "ingress-traefik",
		Tier: namespace.NamespaceTierEdge,
		AdditionalLabels: []namespace.NamespaceLabel{
			{Name: "app.kubernetes.io/name", Value: "ingress-traefik"},
			{Name: "app.kubernetes.io/instance", Value: "ingress-traefik"},
		},
	}
)

func CreateTraefikIngressController(ctx *pulumi.Context) error {
	_, err := namespace.CreateNamespace(ctx, NamespaceTraefikIngressController)
	if err != nil {
		return err
	}

	//// TODO functionalize
	// _, err = createMiddleWares(ctx)
	// if err != nil {
	// 	return err
	// }

	_, err = createTraefikRelease(ctx)
	if err != nil {
		return err
	}

	err = createDashboardIngressRoute(ctx)
	if err != nil {
		return err
	}

	return nil
}

// //TODO basicAuth
func createMiddleWares(ctx *pulumi.Context) ([]*traefikCRD.Middleware, error) {
	middlewares := []*traefikCRD.Middleware{}

	httpRedirect, err := traefikCRD.NewMiddleware(ctx, "httpredirect", &traefikCRD.MiddlewareArgs{
		ApiVersion: pulumi.String("traefik.containo.us/v1alpha1"),
		Kind:       pulumi.String("Middleware"),
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String("redirectscheme"),
			Namespace: pulumi.String(NamespaceTraefikIngressController.Name),
		},
		Spec: &traefikCRD.MiddlewareSpecArgs{
			RedirectScheme: &traefikCRD.MiddlewareSpecRedirectSchemeArgs{
				Permanent: pulumi.Bool(true),
				Scheme:    pulumi.String("https"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	middlewares = append(middlewares, httpRedirect)

	return middlewares, nil
}

func createDashboardIngressRoute(ctx *pulumi.Context) error {
	nameSpaceTraefik := NamespaceTraefikIngressController.Name
	dashboardRouteConf := types.IngressRouteConfig{
		Service: types.Service{
			Kind: types.ServiceKindTraefik,
			Name: "api@internal",
		},
		Name:          "traefik",
		NamespaceName: nameSpaceTraefik,
		Match:         "Host(`traefik" + kubeconfig.DomainNameSuffix(ctx) + "`) && (PathPrefix(`/dashboard`) || PathPrefix(`/api`))",
	}

	_, err := ingressroute.CreateIngressRoute(ctx, &dashboardRouteConf)
	if err != nil {
		return err
	}

	return nil
}

func createTraefikRelease(ctx *pulumi.Context) (*helm.Release, error) {
	traefikRelease, err := helm.NewRelease(ctx, "traefik-ingress", &helm.ReleaseArgs{
		Chart:     pulumi.String("traefik"),
		Namespace: pulumi.String(NamespaceTraefikIngressController.Name),
		Version:   pulumi.String("10.24.3"),
		RepositoryOpts: helm.RepositoryOptsArgs{
			Repo: pulumi.String("https://helm.traefik.io/traefik"),
		},
		Values: pulumi.Map{
			"persistence": pulumi.Map{"enabled": pulumi.Bool(true)},
			"ports": pulumi.Map{
				"web": pulumi.Map{
					"exposedPort": pulumi.Int(ingressControllerPortHTTP),
					//// globalHTTPredirection
					"redirectTo": pulumi.String("redirectsecure"),
					// "redirectTo": pulumi.String("websecure"), //// results in a redirection to 30443 - not available on hostLevel
					// "redirectTo": pulumi.String("443"),       //// lookup Value.ports with index 443 - cannot be found; result:
					//                                           //// POD:spec:containers:-args:-"--entrypoints.web.http.redirections.entryPoint.to=:"
				},
				"websecure": pulumi.Map{
					"exposedPort": pulumi.Int(ingressControllerPortHTTPS),
				},
				"redirectsecure": pulumi.Map{
					"exposedPort": pulumi.Int(443),  //nolint:gomnd
					"port":        pulumi.Int(9999), //nolint:gomnd //// required to build the container
					//// only used as target for https-redirection: the created port is never reached or used.
					//// needed for deployment.spec.template.spec.containers[0].args:-"--entrypoints.web.http.redirections.entryPoint.to=:443"
					//// all traffic hitting 443 gets forwarded to 3{0-2}443 by the hostTraefik and then forwarded to 30443 by k3d
					//// for acmeCerts an clusterissuer's ingressTemplateAnnotation is also needed: traefik.ingress.kubernetes.io/router.tls: "true"
					//// challenges will follow if a tmpCert is available
				},
			},
			"providers": pulumi.Map{
				"kubernetesIngress": pulumi.Map{
					"publishedService": pulumi.Map{
						"enabled": pulumi.Bool(true),
					},
				},
				"kubernetesCRD": pulumi.Map{
					"enabled":             pulumi.Bool(true),
					"allowCrossNamespace": pulumi.Bool(true), //// otherwise the ingressRoute has to be in the same namespace as the service
				},
				"kubernetesGateway": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
			},
			"ingressClass": pulumi.Map{
				"enabled":            pulumi.Bool(true),
				"isDefaultClass":     pulumi.Bool(true),
				"fallbackApiVersion": pulumi.String("v1"),
			},
			"logs": pulumi.Map{
				"general": pulumi.Map{
					// "level": pulumi.String("ERROR"),
					"level": pulumi.String("INFO"),
				},
				"access": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
			},
			"experimental": pulumi.Map{
				//// creats the gatewayclass and an unused gateway
				"kubernetesGateway": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return traefikRelease, nil
}
