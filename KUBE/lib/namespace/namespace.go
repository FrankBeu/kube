// Package lib bundles all types and functions for the pulumi-k8s-definition
package namespace

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NamespaceTier is used to group namespaces respectively apps belonging to different namespaces
// e.g.: show all all available communicationTools
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

	namespaceAPIVersion   string = "v1"
	namespaceKind         string = "Namespace"
	namespaceGlooLabelKey string = "discovery.solo.io/function_discovery"
)

type NamespaceLabel struct {
	Name  string
	Value string
}

// Namespace is used to configure the creation of a pulumiNamespace.
// Remember: bools default to false if not set at instantiation otherwise
type Namespace struct {
	Name             string
	Tier             NamespaceTier
	GlooDiscovery    bool
	AdditionalLabels []NamespaceLabel
}

func CreateNamespace(ctx *pulumi.Context, n *Namespace) (*corev1.Namespace, error) {
	labels := pulumi.StringMap{
		"name": pulumi.String(n.Name),
		"tier": pulumi.String(n.Tier.String()),
	}

	if n.GlooDiscovery {
		labels[namespaceGlooLabelKey] = pulumi.String("enabled")
	}

	if len(n.AdditionalLabels) > 0 {
		for _, label := range n.AdditionalLabels {
			labels[label.Name] = pulumi.String(label.Value)
		}
	}

	ns, err := corev1.NewNamespace(
		ctx,
		n.Name,
		&corev1.NamespaceArgs{
			ApiVersion: pulumi.String(namespaceAPIVersion),
			Kind:       pulumi.String(namespaceKind),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(n.Name),
				Labels: labels,
			},
			Spec: nil,
		},
	)

	if err != nil {
		return nil, err
	}

	return ns, nil
}
