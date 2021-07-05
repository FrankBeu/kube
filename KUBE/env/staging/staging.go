// Package staging contains the configuration for the stagingKube
package staging

import (
	"thesym.site/kube/definition/app/observer/jaeger"
	"thesym.site/kube/definition/app/vcs/gitea"
	"thesym.site/kube/definition/structural/certs/certmanager"
	"thesym.site/kube/definition/structural/ingress/nginx"
	"thesym.site/kube/lib"
	"thesym.site/kube/definition/testing/testnamespace"
	"thesym.site/kube/definition/testing/certtest"
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
	"namestpaceTest": testnamespace.CreateTestNamespace,

	//////////////////////// ////////////////////////
	//// TESTS, PROTOS, ...
	////
	"certTest": certtest.CreateTestCert,

	// "glooPetstore": petstore.CreateGlooPetstore,

	// "fileSingle": pulumiexamples.CreateFromFileSingle,
	// "filesMulti": pulumiexamples.CreateFromFilesMulti,
	// "helmChart": pulumiexamples.CreateFromHelmChart,
	// "fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
