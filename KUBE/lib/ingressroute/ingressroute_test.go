package ingressroute

import (
	"regexp"
	"strconv"
	"sync"
	"testing"

	certv1 "thesym.site/kube/crds/kubernetes/certmanager/v1"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"thesym.site/kube/crds/kubernetes/traefik/v1alpha1"
	"thesym.site/kube/lib/testutil"
	"thesym.site/kube/lib/types"
)

//nolint:funlen
func Test_createIngressRoute(t *testing.T) {
	domain := "domain.test"
	domainNameSuffix := "." + domain
	entrypoint := "websecure"
	type args struct {
		ingRtCnf types.IngressRouteConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test an ingressroute with defaults",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "testingressroute",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
					Middlewares: []types.Middleware{},
				},
			},
		},
		{
			name: "test an ingressroute with dynamic (other) values",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "testingressrouteALT",
					NamespaceName:     "testALT",
					HostSubdomainName: "testingressrouteALT",
					Service: types.Service{
						Name:          "testserviceALT",
						NamespaceName: "testALT",
						Port:          8080,
					},
					Middlewares: []types.Middleware{},
				},
			},
		},
		{
			name: "test an ingressroute with one middleware",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "testingressroute",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
					Middlewares: []types.Middleware{
						{
							Name: types.MiddleWareNameBasicAuth,
						},
					},
				},
			},
		},
		{
			name: "test an ingressroute with multiple middlewares",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "testingressroute",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
					Middlewares: []types.Middleware{
						{
							Name: types.MiddleWareNameBasicAuth,
						},
						{
							Name: types.MiddleWareNameTest,
						},
					},
				},
			},
		},
		{
			name: "test an ingressroute with a custom match",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "testingressroute",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
					Match: "Host(`testingressroute" + domainNameSuffix + "`) && (PathPrefix(`/path1`) || PathPrefix(`/path2`))",
				},
			},
		},
		{
			name: "test an ingressroute w/o custom match and different Name/HostSubdomainName",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "test",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
				},
			},
		},
		{
			name: "test an ingressroute w/  custom match and different Name/HostSubdomainName",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "test",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
					Match: "Host(`testingressroute" + domainNameSuffix + "`) && (PathPrefix(`/path1`) || PathPrefix(`/path2`))",
				},
			},
		},
		{
			name: "test an ingressroute w/o specified ServiceKind",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "test",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Kind:          types.ServiceKindK8S,
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
				},
			},
		},
		{
			name: "test an ingressroute w/ specified ServiceKindK8S",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "test",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Kind:          types.ServiceKindK8S,
						Name:          "testservice",
						NamespaceName: "test",
						Port:          80,
					},
				},
			},
		},
		{
			name: "test an ingressroute w/ specified ServiceKindTraefik",
			args: args{
				ingRtCnf: types.IngressRouteConfig{
					Name:              "test",
					NamespaceName:     "test",
					HostSubdomainName: "testingressroute",
					Service: types.Service{
						Kind: types.ServiceKindTraefik,
						Name: "testservice",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testConfig = map[string]string{
				"project:domain":       ` { "clusterIssuer": "ca-local" }`,
				"project:domainSecret": ` {"domain": "` + domain + `"}`,
			}

			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				cert, ingResult, err := createIngressRoute(ctx, &tt.args.ingRtCnf, domainNameSuffix)
				assert.NoError(t, err)
				var wg sync.WaitGroup
				wg.Add(1)

				//nolint:lll
				// pulumi.All(ingResult.ApiVersion, ingResult.Kind, ingResult.Metadata, ingResult.Spec, cert.Metadata).ApplyT(func(all []interface{}) error {
				pulumi.All(ingResult.ApiVersion, ingResult.Kind, ingResult.Metadata, ingResult.Spec, cert.Spec).ApplyT(func(all []interface{}) error {
					apivActual := *all[0].(*string)
					kindActual := *all[1].(*string)
					metaActual := all[2].(metav1.ObjectMeta)
					specActual := all[3].(v1alpha1.IngressRouteSpec)
					certSpecActual := all[4].(certv1.CertificateSpec)

					//nolint:gocritic
					// spew.Dump(apivActual)
					// spew.Dump(kindActual)
					// spew.Dump(specActual)
					// spew.Dump(metaActual)
					// spew.Dump(certSpecActual)

					////staticProperties
					testutil.Equalf(t, "IngressRoute", "ApiVersion", apivActual, types.IngressRouteAPIVersion)
					testutil.Equalf(t, "IngressRoute", "Kind", kindActual, types.IngressRouteKind)
					testutil.Equalf(t, "IngressRoute", "Entrypoint", specActual.EntryPoints[0], entrypoint)

					////dynamicProperties
					testutil.Equalf(t, "IngressRoute", "Name", *metaActual.Name, tt.args.ingRtCnf.Name)
					testutil.Equalf(t, "IngressRoute", "Namespace", *metaActual.Namespace, tt.args.ingRtCnf.NamespaceName)
					testutil.Equalf(t, "IngressRoute", "TLSSecretName", *specActual.Tls.SecretName, tt.args.ingRtCnf.Name+types.IngressRouteTLSSecretSuffix) //nolint:lll

					assertMiddlewares(t, specActual.Routes[0].Middlewares, tt.args.ingRtCnf.Middlewares)

					////ensure ingressroute.tls.secretName == certificate.metadata.name
					assert.Equalf(t, *specActual.Tls.SecretName, certSpecActual.SecretName, "%s's TLSSecretName != CertSpecSecretName: %q -- %q", "IngressRoute", *specActual.Tls.SecretName, certSpecActual.SecretName)

					////if match is set; use it instead
					hostTarget := ""
					if tt.args.ingRtCnf.Match != "" {
						hostTarget = tt.args.ingRtCnf.Match
					} else {
						hostTarget = "Host(`" + tt.args.ingRtCnf.HostSubdomainName + domainNameSuffix + "`)"
					}
					testutil.Equalf(t, "IngressRoute", "Host", specActual.Routes[0].Match, hostTarget) //nolint:lll
					////ensure ingressroute...match == cert.commonName
					hostActual := specActual.Routes[0].Match
					re := regexp.MustCompile(`^Host\(.(.*?).\)`)
					matchClean := re.FindStringSubmatch(hostActual)[1]
					assert.Equalf(t, matchClean, *certSpecActual.CommonName, "%s's Match != CertSpecCommonName: %q -- %q", "IngressRoute", matchClean, *certSpecActual.CommonName)

					assertService(t, specActual.Routes[0].Services, tt.args.ingRtCnf.Service)

					wg.Done()
					return nil
				})

				return nil
			}, testutil.WithMocksAndConfig("project", "stack", testConfig, testutil.Mocks(0)))
			assert.NoError(t, err)
		})
	}
}

func assertMiddlewares(t *testing.T, middlewares []v1alpha1.IngressRouteSpecRoutesMiddlewares, middlewareTargets []types.Middleware) {
	//// generate an array of middlewareNames to compare the middlewareNameTarget (string) with
	var middlewareNamesActual []string
	for _, middleware := range middlewares {
		middlewareNamesActual = append(middlewareNamesActual, middleware.Name)
		//// check Namespace
		//nolint:lll
		// testutil.Equalf(t, "IngressRoute", "MiddlewareNamespace", *middleware.Namespace, traefik.NamespaceTraefikIngressController.Name)////TODO cannot import traefik
		testutil.Equalf(t, "IngressRoute", "MiddlewareNamespace", *middleware.Namespace, "ingress-traefik")
	}

	//// check if middlewareTarget actually exists
	for _, mt := range middlewareTargets {
		testutil.Containsf(t, "IngressRoute", "MiddlewareNames", middlewareNamesActual, mt.Name.String())
	}

	//// check if no other middleware exists
	testutil.Equalf(t, "IngressRoute", "quantity", len(middlewareNamesActual), len(middlewareTargets))
}

func assertService(t *testing.T, services []v1alpha1.IngressRouteSpecRoutesServices, serviceTarget types.Service) {
	//// currently only one service is used
	var serviceActual = services[0]

	testutil.Equalf(t, "IngressRoute", "ServiceName", serviceActual.Name, serviceTarget.Name)
	switch serviceTarget.Kind {
	case types.ServiceKindK8S:
		testutil.Equalf(t, "IngressRoute", "ServiceNamespace", *serviceActual.Namespace, serviceTarget.NamespaceName)
		testutil.Equalf(t, "IngressRoute", "ServicePort", strconv.Itoa(int(serviceActual.Port.(float64))), strconv.Itoa(serviceTarget.Port))

	case types.ServiceKindTraefik:
		assert.Equal(t, interface{}(nil), serviceActual.Port, "IngressRoute has wrong ServicePort: ServiceKindTraefik must not specify a port.")
		assert.Equal(t, (*string)(nil), serviceActual.Namespace,
			"IngressRoute has wrong ServiceNamespace: ServiceKindTraefik must not specify a namespace.")
	}
}
