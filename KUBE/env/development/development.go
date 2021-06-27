// Package development contains the configuration for the developmentKube
package development

import (
	"thesym.site/k8s/definition/structural/certs/certmanager"
	"thesym.site/k8s/definition/structural/ingress/nginx"
	"thesym.site/k8s/definition/testing/gloo/petstore"
	"thesym.site/k8s/lib"
)

// Kube is the configuration for the developmentEnvironment
// var Kube = map[string]func(*pulumi.Context) error{
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
	//// VCS
	////
	// "gitea": gitea.CreateGitea,

	//////////////////////// //////////////////////// ////////////////////////
	//// TESTING
	////
	"glooPetstore": petstore.CreateGlooPetstore,

	// "fileSingle": pulumiexamples.CreateFromFileSingle,
	// "filesMulti": pulumiexamples.CreateFromFilesMulti,
	// "helmChart": pulumiexamples.CreateFromHelmChart,
	// "fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
