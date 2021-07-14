// Package ingress is responsible for creation of ingress-related resources
package ingress

import (
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	cert "thesym.site/kube/lib/certificate"
	"thesym.site/kube/lib/config"
)

//nolint:golint
// IngressClassName specifies the ingressController which implements the ingress
type IngressClassName int

//go:generate stringer -type=IngressClassName -linecomment
const (
	IngressClassNameNginx IngressClassName = iota // nginx
	// IngressClassNameTest // test

	ingressAPIVersion string = "networking.k8s.io/v1"
	ingressKind       string = "Ingress"
	tlsAnnotationKey  string = "cert-manager.io/cluster-issuer"
	tlsSecretSuffix   string = "-tls"
)

type Host struct {
	Name        string
	ServiceName string
	ServicePort int
}

type Config struct {
	Annotations       pulumi.StringMap
	ClusterIssuerType cert.ClusterIssuerType
	Hosts             []Host
	IngressClassName  IngressClassName
	Name              string
	NamespaceName     string
	TLS               bool
}

func CreateIngress(ctx *pulumi.Context, ing *Config) (*networkingv1.Ingress, error) {
	domainNameSuffix := config.DomainNameSuffix(ctx)

	//// TLS annotation
	annotations := pulumi.StringMap{}
	if ing.TLS {
		annotations["cert-manager.io/cluster-issuer"] = pulumi.String(cert.ClusterIssuerTypeCaLocal.String())
	}
	for k, v := range ing.Annotations {
		annotations[k] = v
	}

	tls := networkingv1.IngressTLSArray{}
	tlsHosts := pulumi.StringArray{}
	for _, host := range ing.Hosts {
		tlsHosts = append(tlsHosts, pulumi.String(host.Name+domainNameSuffix))
	}
	if ing.TLS {
		tls = networkingv1.IngressTLSArray{
			&networkingv1.IngressTLSArgs{
				Hosts:      tlsHosts,
				SecretName: pulumi.String(ing.Name + tlsSecretSuffix),
			},
		}
	}

	rules := networkingv1.IngressRuleArray{}
	for _, host := range ing.Hosts {
		hostRule := &networkingv1.IngressRuleArgs{
			Host: pulumi.String(host.Name + domainNameSuffix),
			Http: &networkingv1.HTTPIngressRuleValueArgs{
				Paths: networkingv1.HTTPIngressPathArray{
					&networkingv1.HTTPIngressPathArgs{
						PathType: pulumi.String("Prefix"),
						Path:     pulumi.String("/"),
						Backend: &networkingv1.IngressBackendArgs{
							Service: &networkingv1.IngressServiceBackendArgs{
								Name: pulumi.String(host.ServiceName),
								Port: &networkingv1.ServiceBackendPortArgs{
									Number: pulumi.Int(host.ServicePort),
								},
							},
						},
					},
				},
			},
		}

		rules = append(rules, hostRule)
	}

	ingress, err := networkingv1.NewIngress(ctx, ing.Name, &networkingv1.IngressArgs{
		ApiVersion: pulumi.String(ingressAPIVersion),
		Kind:       pulumi.String(ingressKind),
		Metadata: &metav1.ObjectMetaArgs{
			Labels:    pulumi.StringMap{"app": pulumi.String(ing.Name)},
			Name:      pulumi.String(ing.Name),
			Namespace: pulumi.String(ing.NamespaceName),
			// Annotations: ing.Annotations,
			Annotations: annotations,
		},
		Spec: &networkingv1.IngressSpecArgs{
			IngressClassName: pulumi.String(ing.IngressClassName.String()),
			Rules:            rules,
			Tls:              tls,
		},
	})
	if err != nil {
		return nil, err
	}
	return ingress, nil
}
