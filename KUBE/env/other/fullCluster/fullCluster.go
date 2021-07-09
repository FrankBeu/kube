// Package fullcluster contains the configuration for a cluster with all resources activated.
// Creation of all resources in the cluster will fail because some use the same ports.
// But yaml-files will be generated
package fullcluster

import (
	"thesym.site/kube/definition/app/communication/jitsi"
	"thesym.site/kube/definition/app/communication/matrix"
	"thesym.site/kube/definition/app/observer/jaeger"
	"thesym.site/kube/definition/app/vcs/gitea"
	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/emmissary"
	"thesym.site/kube/definition/structural/ingress/nginx"
	"thesym.site/kube/definition/structural/ingress/tyk"
	"thesym.site/kube/definition/structural/monitoring/prometheus"
	"thesym.site/kube/definition/testing/gloo/petstore"
	"thesym.site/kube/definition/testing/pulumiexamples"
	"thesym.site/kube/definition/testing/testcert"
	"thesym.site/kube/definition/testing/testingress"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/lib"
)

// Kube is the configuration for the stagingEnvironment
var Kube = lib.KubeConfig{
	//////////////////////// //////////////////////// ////////////////////////
	//// STRUCTURAL
	////
	//////////////////////// ////////////////////////
	//// Ingress
	////
	"nginxIngress": nginx.CreateNginxIngressController,

	//// installation working; crd-usage not
	"emmissary": emmissary.CreateEmmissary,

	// NOTREADY gloo

	"tykGateway": tyk.CreateTykGateway,

	//////////////////////// ////////////////////////
	//// certificates
	////
	"certmanager": certmanager.CreateCertmanager,

	//////////////////////// ////////////////////////
	//// MONITORING
	////
	"prometheus": prometheus.CreatePrometheus,

	//////////////////////// //////////////////////// ////////////////////////
	//// APPS
	////

	//////////////////////// ////////////////////////
	//// COMMUNICATION
	////
	"jitsi": jitsi.CreateJitsi,

	////////////////////////
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
	"testNamespace": testnamespace.CreateTestNamespace,

	//////////////////////// ////////////////////////
	//// TESTS, PROTOS, ...
	////
	"testCert": testcert.CreateTestCert,
	"testIngress": testingress.CreateTestIngress,

	"glooPetstore": petstore.CreateGlooPetstore,

	"fileSingle":   pulumiexamples.CreateFromFileSingle,
	"filesMulti":   pulumiexamples.CreateFromFilesMulti,
	"helmChart":    pulumiexamples.CreateFromHelmChart,
	"fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
