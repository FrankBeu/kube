// Package ingress is responsible for creation of ingress-related resources
package types

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//nolint:golint
// IngressClassName specifies the ingressController which implements the ingress
type IngressClassName int

//go:generate stringer -type=IngressClassName -linecomment
const (
	IngressClassNameNginx IngressClassName = iota // nginx
	// IngressClassNameTest // test

	IngressAPIVersion string = "networking.k8s.io/v1"
	IngressKind       string = "Ingress"
	TlsAnnotationKey  string = "cert-manager.io/cluster-issuer"
	TlsSecretSuffix   string = "-tls"
)

type Host struct {
	Name        string
	ServiceName string
	ServicePort int
}

type Config struct {
	Annotations       pulumi.StringMap
	ClusterIssuerType ClusterIssuerType
	Hosts             []Host
	IngressClassName  IngressClassName
	Name              string
	NamespaceName     string
	TLS               bool
}
