// Package development contains the configuration for the developmentKube
package development

import (
	// "thesym.site/kube/definition/app/communication/matrix"
	// "thesym.site/kube/definition/structural/monitoring/prometheus"
	// "thesym.site/kube/definition/testing/testnamespace"

	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/traefik"

	// "thesym.site/kube/definition/testing/testingress"

	"thesym.site/kube/definition/testing/testcert"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/definition/testing/whoami"
	"thesym.site/kube/lib/kubeConfig"
)

// Kube is the configuration for the developmentEnvironment
// var Kube = map[string]func(*pulumi.Context) error{
var Kube = kubeConfig.KubeConfig{
	//////////////////////// //////////////////////// ////////////////////////
	//// STRUCTURAL
	////
	//////////////////////// ////////////////////////
	//// Ingress
	////
	"traefikIngress": traefik.CreateTraefikIngressController,

	//////////////////////// ////////////////////////
	//// certificates
	////
	"certmanager": certmanager.CreateCertmanager,

	//////////////////////// ////////////////////////
	//// MONITORING
	////
	// "prometheus": prometheus.CreatePrometheus,

	//////////////////////// //////////////////////// ////////////////////////
	//// APPS
	////

	//////////////////////// ////////////////////////
	//// COMMUNICATION
	////
	// "jitsi": jitsi.CreateJitsi,
	// "matrix": matrix.CreateMatrix,

	////////////////////////
	//// OBSERVING
	////
	// "jaeger": jaeger.CreateJaegerOperator,

	//////////////////////// ////////////////////////
	//// VCS
	////
	// "gitea": gitea.CreateGitea,

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
	"testCert": testcert.CreateTestCert,
	// "testIngress": testingress.CreateTestIngress,
	// "testHelmRelease": testhelmrelease.CreateHelmRelease,
	"whoami": whoami.CreateWhoAmI,

	// "glooPetstore": petstore.CreateGlooPetstore,

	// "fileSingle": pulumiexamples.CreateFromFileSingle,
	// "filesMulti": pulumiexamples.CreateFromFilesMulti,
	// "helmChart": pulumiexamples.CreateFromHelmChart,
	// "fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
