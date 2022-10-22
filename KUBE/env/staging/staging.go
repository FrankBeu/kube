// Package staging contains the configuration for the stagingKube
package staging

import (
	// "thesym.site/kube/definition/app/communication/jitsi"
	// "thesym.site/kube/definition/app/communication/matrix"
	"thesym.site/kube/definition/app/observer/jaeger"
	"thesym.site/kube/definition/app/vcs/gitea"
	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/traefik"

	// "thesym.site/kube/definition/structural/monitoring/prometheus"
	// "thesym.site/kube/definition/testing/gloo/petstore"
	// "thesym.site/kube/definition/testing/testingress"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/definition/testing/whoami"
	"thesym.site/kube/lib/kubeConfig"
)

// Kube is the configuration for the stagingEnvironment
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
	// "loki": loki.CreateLoki,

	//////////////////////// //////////////////////// ////////////////////////
	//// APPS
	////

	//////////////////////// ////////////////////////
	//// COMMUNICATION
	////
	// "jitsi": jitsi.CreateJitsi,
	// "matrix": matrix.CreateMatrix,

	//////////////////////// ////////////////////////
	//// OBSERVING
	////
	"jaeger": jaeger.CreateJaegerOperator,

	//////////////////////// ////////////////////////
	//// VCS
	////
	"gitea": gitea.CreateGitea,

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
	"whoami": whoami.CreateWhoAmI,

	// "glooPetstore": petstore.CreateGlooPetstore,

	// "fileSingle": pulumiexamples.CreateFromFileSingle,
	// "filesMulti": pulumiexamples.CreateFromFilesMulti,
	// "helmChart": pulumiexamples.CreateFromHelmChart,
	// "fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
