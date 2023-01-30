// TODO ingress is not currently not working with traefik as ingress controller
package types

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IngressClassName specifies the ingressController which implements the ingress
type IngressClassName int
type IngressProperty string

//go:generate stringer -type=IngressClassName -linecomment
const (
	IngressClassNameTraefik IngressClassName = iota // traefik
	IngressClassNameNginx                           // nginx
	// IngressClassNameTest                         // test

	IngressAPIVersion       string = "networking.k8s.io/v1"
	IngressKind             string = "Ingress"
	IngressTLSAnnotationKey string = "cert-manager.io/cluster-issuer"
	IngressTLSSecretSuffix  string = "-tls"
)

type Host struct {
	Name        string // used as DomainName
	ServiceName string
	ServicePort int
}

type IngressConfig struct {
	Annotations       pulumi.StringMap
	ClusterIssuerType ClusterIssuerType
	Hosts             []Host
	IngressClassName  IngressClassName
	Name              string
	NamespaceName     string
	TLS               bool
}
