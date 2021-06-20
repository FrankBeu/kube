// Package lib bundles all types and functions for the pulumi-k8s-definition
package lib

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NamespaceTier is used to group namespaces respectively apps belonging to different namespaces
type NamespaceTier int

// onChange: regenerate with `go generate`
//go:generate stringer -type=NamespaceTier -linecomment
const (
	NamespaceTierAPI               NamespaceTier = iota // api
	NamespaceTierApp                                    // app
	NamespaceTierAuth                                   // auth
	NamespaceTierCommunication                          // communication
	NamespaceTierDatabase                               // database
	NamespaceTierDevOps                                 // devOps
	NamespaceTierDNS                                    // dns
	NamespaceTierEdge                                   // edge
	NamespaceTierMonitoring                             // monitoring
	NamespaceTierProjectManagement                      // projectManagement
	NamespaceTierSecurity                               // security
	NamespaceTierStorage                                // storage
	NamespaceTierTesting                                // testing
)

type Namespace struct {
	Name string
	Tier NamespaceTier
}

func CreateNamespace(n *Namespace, ctx *p.Context) error {
	labels := p.StringMap{
		"name": p.String(n.Name),
		"tier": p.String(n.Tier.String()),
	}

	_, err := corev1.NewNamespace(
		ctx,
		n.Name,
		&corev1.NamespaceArgs{
			ApiVersion: p.String("v1"),
			Kind:       p.String("Namespace"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   p.String(n.Name),
				Labels: labels,
			},
			Spec: nil,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
