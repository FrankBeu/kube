package ingress

import (
	"strconv"
	"sync"
	"testing"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	cert "thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/config"
	"thesym.site/kube/lib/testutil"
)

//nolint:funlen
func TestCreateIngress(t *testing.T) {
	type args struct {
		ing Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test a ingress with defaults",
			args: args{
				ing: Config{
					Name:             "testingress",
					NamespaceName:    "test",
					IngressClassName: IngressClassNameNginx,
					Hosts: []Host{
						{
							Name:        "test",
							ServiceName: "test",
							ServicePort: 8888,
						},
					},
				},
			},
		},
		{
			name: "test an ingress with two hosts",
			args: args{
				ing: Config{
					Name:             "otheringress",
					NamespaceName:    "other",
					IngressClassName: IngressClassNameNginx,

					Hosts: []Host{
						{
							Name:        "test",
							ServiceName: "test",
							ServicePort: 8888,
						},
						{
							Name:        "other",
							ServiceName: "other",
							ServicePort: 9999,
						},
					},
				},
			},
		},
		{
			name: "test an ingress with multiple hosts",
			args: args{
				ing: Config{
					Name:             "another",
					NamespaceName:    "test",
					IngressClassName: IngressClassNameNginx,
					Hosts: []Host{
						{
							Name:        "test",
							ServiceName: "test",
							ServicePort: 8888,
						},
						{
							Name:        "other",
							ServiceName: "other",
							ServicePort: 9999,
						},
						{
							Name:        "another",
							ServiceName: "another",
							ServicePort: 7777,
						},
					},
				},
			},
		},
		{
			name: "test an ingress with annotations",
			args: args{
				ing: Config{
					Name:             "annotationingress",
					NamespaceName:    "other",
					IngressClassName: IngressClassNameNginx,
					Annotations: pulumi.StringMap{
						"nginx.ingress.kubernetes.io/ssl-redirect":       pulumi.String("false"),
						"nginx.ingress.kubernetes.io/force-ssl-redirect": pulumi.String("false"),
					},
					Hosts: []Host{
						{
							Name:        "test",
							ServiceName: "test",
							ServicePort: 8888,
						},
					},
				},
			},
		},
		{
			name: "test an ingress with tls activated",
			args: args{
				ing: Config{
					Name:             "tlsingress",
					NamespaceName:    "other",
					IngressClassName: IngressClassNameNginx,
					Hosts: []Host{
						{
							Name:        "test",
							ServiceName: "test",
							ServicePort: 8888,
						},
					},
					TLS: true,
				},
			},
		},
		{
			name: "test an ingress with tls activated and additional annotations",
			args: args{
				ing: Config{
					Name:             "otheringress",
					NamespaceName:    "other",
					IngressClassName: IngressClassNameNginx,
					Annotations: pulumi.StringMap{
						"nginx.ingress.kubernetes.io/force-ssl-redirect": pulumi.String("true"),
					},
					Hosts: []Host{
						{
							Name:        "test",
							ServiceName: "test",
							ServicePort: 8888,
						},
					},
					TLS: true,
				},
			},
		},
		{
			name: "test an ingress with tls activated and multiple hosts",
			args: args{
				ing: Config{
					Name:             "othertlsingress",
					NamespaceName:    "other",
					IngressClassName: IngressClassNameNginx,
					Hosts: []Host{
						{
							Name:        "test",
							ServiceName: "test",
							ServicePort: 8888,
						},
						{
							Name:        "other",
							ServiceName: "other",
							ServicePort: 9999,
						},
					},
					TLS: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				ingResult, err := CreateIngress(ctx, &tt.args.ing)
				assert.NoError(t, err)

				var wg sync.WaitGroup
				wg.Add(1)

				pulumi.All(ingResult.ApiVersion, ingResult.Kind, ingResult.Metadata, ingResult.Spec).ApplyT(func(all []interface{}) error {
					apivActual := *all[0].(*string)
					kindActual := *all[1].(*string)
					metaActual := *all[2].(*metav1.ObjectMeta)
					specActual := *all[3].(*networkingv1.IngressSpec)

					testutil.Equalf(t, "Ingress", "ApiVersion", apivActual, ingressAPIVersion)
					testutil.Equalf(t, "Ingress", "Kind", kindActual, ingressKind)
					testutil.Equalf(t, "Ingress", "Name", *metaActual.Name, tt.args.ing.Name)
					testutil.Equalf(t, "Ingress", "Namespace", *metaActual.Namespace, tt.args.ing.NamespaceName)
					testutil.Equalf(t, "Ingress", "IngressClassName", *specActual.IngressClassName, tt.args.ing.IngressClassName.String())

					assertSpecRules(t, ctx, specActual.Rules, tt.args.ing.Hosts)

					assertAnnotations(t, tt.args.ing.TLS, metaActual.Annotations, tt.args.ing.Annotations)

					if tt.args.ing.TLS {
						assertSpecTLS(t, ctx, specActual.Tls, tt.args.ing.Hosts, tt.args.ing.Name+tlsSecretSuffix)
					} else {
						assert.Emptyf(t, specActual.Tls, "TLS: specTls is not empty despite tls NOT being activated: actual: %+v", specActual.Tls)
					}

					wg.Done()
					return nil
				})

				return nil
			}, testutil.WithMocksAndConfig("project", "stack", testutil.TestConfig, testutil.Mocks(0)))
			assert.NoError(t, err)
		})
	}
}

func assertSpecRules(t *testing.T, ctx *pulumi.Context, rules []networkingv1.IngressRule, hostsTarget []Host) {
	//// extract the necessary data
	var hostsActual []Host
	for _, rule := range rules {
		hostActual := Host{
			Name: *rule.Host,
			//// currently only Hosts with a single path are used
			ServiceName: rule.Http.Paths[0].Backend.Service.Name,
			ServicePort: *rule.Http.Paths[0].Backend.Service.Port.Number,
		}
		hostsActual = append(hostsActual, hostActual)
	}
	var hostNamesActual []string
	for _, rule := range rules {
		hostNamesActual = append(hostNamesActual, *rule.Host)
	}

	//// compare
	for _, ht := range hostsTarget {
		hostDomainNameTarget := ht.Name + config.DomainNameSuffix(ctx)

		testutil.Containsf(t, "Ingress", "Rules:Host", hostNamesActual, hostDomainNameTarget)
		for _, ha := range hostsActual {
			if ha.Name == hostDomainNameTarget {
				testutil.Equalf(t, "ingress", "rules:http:service:name", ha.ServiceName, ht.ServiceName)
				testutil.Equalf(t, "ingress", "rules:http:service:port", strconv.Itoa(ha.ServicePort), strconv.Itoa(ht.ServicePort))
			}
		}
	}
}

func assertAnnotations(t *testing.T, tlsEnabled bool, annotationsActual map[string]string, annotationsTarget pulumi.StringMap) {
	if tlsEnabled {
		containsMsg := "TLS: annotation %q is not set"
		testutil.AssertAnnotation(t, annotationsActual, tlsAnnotationKey, pulumi.String(cert.ClusterIssuerTypeCaLocal.String()), containsMsg)
	} else {
		assert.NotContainsf(
			t,
			annotationsActual,
			tlsAnnotationKey,
			"TLS: label %q is set despite TLS not being activated",
			tlsAnnotationKey,
		)
	}

	containsMsg := "annotation %q is not set"
	for a, v := range annotationsTarget {
		testutil.AssertAnnotation(t, annotationsActual, a, v, containsMsg)
	}
}

func assertSpecTLS(t *testing.T, ctx *pulumi.Context, tlsActual []networkingv1.IngressTLS, hostsTarget []Host, secretName string) {
	for _, ht := range hostsTarget {
		hostDomainNameTarget := ht.Name + config.DomainNameSuffix(ctx)

		//// currently only one TLS-item is used
		testutil.Containsf(t, "Ingress", "Tls[0].Hosts", tlsActual[0].Hosts, hostDomainNameTarget)
	}
	testutil.Equalf(t, "Ingress", "tls.SecretName", *tlsActual[0].SecretName, secretName)
}
