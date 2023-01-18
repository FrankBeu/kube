// Package ingress is responsible for creation of ingress-related resources
package types

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IngressClassName specifies the ingressController which implements the ingress
type IngressClassName int

//go:generate stringer -type=IngressClassName -linecomment
const (
	IngressClassNameTraefik IngressClassName = iota //// traefik
	IngressClassNameNginx                           //// nginx
	// IngressClassNameTest                         //// test

	IngressAPIVersion string = "networking.k8s.io/v1"
	IngressKind       string = "Ingress"
	TLSAnnotationKey  string = "cert-manager.io/cluster-issuer"
	TLSSecretSuffix   string = "-tls"
)

type Host struct {
	Name        string
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
