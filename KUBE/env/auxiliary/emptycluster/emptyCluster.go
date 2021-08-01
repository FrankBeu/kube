// Package emptycluster is used to remove all resources from a cluster
package emptycluster

import (
	"thesym.site/kube/lib/kubeConfig"
)

// Kube is the configuration for an empty cluster
var Kube = kubeConfig.KubeConfig{}
