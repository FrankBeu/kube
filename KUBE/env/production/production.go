// Package production contains the configuration for the productionKube
package production

import (
	"thesym.site/kube/definition/app/vcs/gitea"
	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/nginx"
	"thesym.site/kube/lib"
)

// Kube is the configuration for the productionEnvironment
var Kube = lib.KubeConfig{
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

	//////////////////////// //////////////////////// ////////////////////////
	//// APPS
	////

	//////////////////////// ////////////////////////
	//// COMMUNICATION
	////
	// "jitsi": jitsi.CreateJitsi,

	//////////////////////// ////////////////////////
	//// MONITORING
	////
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
	// "testNamespace": testnamespace.CreateTestNamespace,

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
