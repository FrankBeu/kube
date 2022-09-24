// Package development contains the configuration for the developmentKube
package development

import (
	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/nginx"

	"thesym.site/kube/definition/testing/testingress"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/lib/kubeConfig"
)

// Kube is the configuration for the developmentEnvironment
var Kube = kubeConfig.KubeConfig{
	//////////////////////// //////////////////////// ////////////////////////
	//// STRUCTURAL
	////
	//////////////////////// ////////////////////////
	//// Ingress
	////
	"nginxIngress": nginx.CreateNginxIngressController,

	//////////////////////// ////////////////////////
	//// certificates
	////
	"certmanager": certmanager.CreateCertmanager,

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
	"testIngress": testingress.CreateTestIngress,
}
