// Package development contains the configuration for the developmentKube
package development

import (
	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/traefik"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/definition/testing/whoami"
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

	//////////////////////// ////////////////////////
	//// CERTIFICATES
	////
	"certmanager": certmanager.CreateCertmanager,

	//////////////////////// ////////////////////////
	//// MONITORING
	////

	//////////////////////// ////////////////////////
	//// OBSERVING
	////

	//////////////////////// //////////////////////// ////////////////////////
	//// APPS
	////

	//////////////////////// ////////////////////////
	//// COMMUNICATION
	////

	//////////////////////// ////////////////////////
	//// VCS
	////

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
	// "testCert": testcert.CreateTestCert,
	// "testIngress": testingress.CreateTestIngress,
	// "testHelmRelease": testhelmrelease.CreateHelmRelease,

	"whoami":             whoami.CreateWhoAmI, //// whoami does not create a way to ingress; choose one:
	"whoamiGateway":      whoami.CreateGateway,
	"whoamiIngressRoute": whoami.CreateIngressRoute,

	// "fileSingle": pulumiexamples.CreateFromFileSingle,
	// "filesMulti": pulumiexamples.CreateFromFilesMulti,
	// "helmChart": pulumiexamples.CreateFromHelmChart,
	// "fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
