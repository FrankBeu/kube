// Package emptycluster is used to remove all resources from a cluster
package emptycluster

import (
	"thesym.site/kube/lib/kubeconfig"
)

// Kube is the configuration for an empty cluster
var Kube = kubeconfig.KubeConfig{}
