// Package development contains the configuration for the developmentKube
package development

import (
	"thesym.site/kube/definition/structural/ingress/traefik"
	"thesym.site/kube/definition/testing/pulumiexamples"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/lib/kubeconfig"
)

// Kube is the configuration for the developmentEnvironment
var Kube = kubeconfig.KubeConfig{
	//////////////////////// //////////////////////// ////////////////////////
	//// STRUCTURAL
	////
	//////////////////////// ////////////////////////
	//// INGRESS
	////
	"traefikIngress": traefik.CreateTraefikIngressController,

	//////////////////////// //////////////////////// ////////////////////////
	//// TESTING
	////
	//////////////////////// ////////////////////////
	//// NAMESPACE
	////
	"testNamestpace": testnamespace.CreateTestNamespace,

	//////////////////////// ////////////////////////
	//// TESTS, PROTOS, ...
	////
	"filesMulti": pulumiexamples.CreateFromFilesMulti,
}
