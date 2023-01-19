// Package fullcluster contains the configuration for a cluster with all resources activated.
// Creation of all resources in the cluster will fail because some use the same ports.
// But yaml-files will be generated
// TODO reenable
package fullcluster

import (
	"thesym.site/kube/lib/kubeConfig"
	// "thesym.site/kube/definition/app/communication/jitsi"
	// "thesym.site/kube/definition/app/communication/matrix"
	// "thesym.site/kube/definition/app/communication/nextcloud"
	// "thesym.site/kube/definition/app/communication/owncloud"
	// "thesym.site/kube/definition/app/observer/jaeger"
	// "thesym.site/kube/definition/app/vcs/gitea"
	// "thesym.site/kube/definition/structural/certs/certmanager"
	// "thesym.site/kube/definition/structural/ingress/traefik"
	// "thesym.site/kube/definition/structural/monitoring/loki"
	// "thesym.site/kube/definition/structural/monitoring/prometheus"
	// "thesym.site/kube/definition/testing/pulumiexamples"
	// "thesym.site/kube/definition/testing/testcert"
	// "thesym.site/kube/definition/testing/testhelmrelease"
	// "thesym.site/kube/definition/testing/testingress"
	// "thesym.site/kube/definition/testing/testnamespace"
	// "thesym.site/kube/definition/testing/whoami"
	// "thesym.site/kube/lib/kubeConfig"
)

// Kube is the configuration for the stagingEnvironment
var Kube = kubeConfig.KubeConfig{
	//////////////////////// //////////////////////// ////////////////////////
	//// STRUCTURAL
	////
	//////////////////////// ////////////////////////
	//// INGRESS
	////
	// "traefikIngress": traefik.CreateTraefikIngressController,

	//////////////////////// ////////////////////////
	//// CERTIFICATES
	////
	// "certmanager": certmanager.CreateCertmanager,

	//////////////////////// ////////////////////////
	//// MONITORING
	////
	// "prometheus": prometheus.CreatePrometheus,
	// "loki":       loki.CreateLoki,

	//////////////////////// ////////////////////////
	//// OBSERVING
	////
	// "jaeger": jaeger.CreateJaegerOperator,

	//////////////////////// //////////////////////// ////////////////////////
	//// APPS
	////

	//////////////////////// ////////////////////////
	//// COMMUNICATION
	////
	// "jitsi":  jitsi.CreateJitsi,
	// "matrix": matrix.CreateMatrix,
	// "nextCloud": nextcloud.CreateNextCloud,
	// "ownCloud":  owncloud.CreateOwnCloud,

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
	// "testNamespace": testnamespace.CreateTestNamespace,
	//////////////////////// ////////////////////////
	//// TESTS, PROTOS, ...
	////
	// "testCert":        testcert.CreateTestCert,
	// "testHelmRelease": testhelmrelease.CreateHelmRelease,
	// "testIngress":     testingress.CreateTestIngress,
	// "whoami":          whoami.CreateWhoAmI,

	// "fileSingle":   pulumiexamples.CreateFromFileSingle,
	// "filesMulti":   pulumiexamples.CreateFromFilesMulti,
	// "helmChart":    pulumiexamples.CreateFromHelmChart,
	// "fileSingleWT": pulumiexamples.CreateFromFileSingleWithConfigurationTransformation,
}
