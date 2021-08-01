// Package development contains the configuration for the developmentKube
package development

import (
	// "thesym.site/kube/definition/app/communication/matrix"
	// "thesym.site/kube/definition/structural/monitoring/prometheus"
	// "thesym.site/kube/definition/testing/testnamespace"

	"thesym.site/kube/definition/app/vcs/gitea"
	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/nginx"

	// "thesym.site/kube/definition/testing/testingress"
	"thesym.site/kube/definition/testing/testnamespace"
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
	"nginxIngress": nginx.CreateNginxIngressController,

	//// installation working; crd-usage not
	// "emmissary": emmissary.CreateEmmissary,

	// NOTREADY gloo

	// "tykGateway": tyk.CreateTykGateway,

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

	// "glooPetstore": petstore.CreateGlooPetstore,

	// "fileSingle": pulumiexamples.CreateFromFileSingle,
	// "filesMulti": pulumiexamples.CreateFromFilesMulti,
	// "helmChart": pulumiexamples.CreateFromHelmChart,
	// "fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
