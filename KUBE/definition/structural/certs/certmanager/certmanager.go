// Package certmanager provides certmanager-crds and the default certmanager-installation
package certmanager

import (
	"fmt"

	admissionregistrationv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/admissionregistration/v1"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	certv1b1 "thesym.site/kube/crds/cert-manager/certmanager/v1beta1"
	"thesym.site/kube/lib"
)

func CreateCertmanager(ctx *pulumi.Context) error {
	//nolint:lll
	certManagerCrds := []lib.Crd{
		{Name: "certmanager-order-definition", Location: "./crds/cert-manager/cdrDefinitions/customresourcedefinition-orders.acme.cert-manager.io.yaml"},
		{Name: "certmanager-issuer-definition", Location: "./crds/cert-manager/cdrDefinitions/customresourcedefinition-issuers.cert-manager.io.yaml"},
		{Name: "certmanager-challenge-definition", Location: "./crds/cert-manager/cdrDefinitions/customresourcedefinition-challenges.acme.cert-manager.io.yaml"},
		{Name: "certmanager-certificate-definition", Location: "./crds/cert-manager/cdrDefinitions/customresourcedefinition-certificates.cert-manager.io.yaml"},
		{Name: "certmanager-clusterIssuer-definition", Location: "./crds/cert-manager/cdrDefinitions/customresourcedefinition-clusterissuers.cert-manager.io.yaml"},
		{Name: "certmanager-certificate-request-definition", Location: "./crds/cert-manager/cdrDefinitions/customresourcedefinition-certificaterequests.cert-manager.io.yaml"},
	}
	//// Register the CRD.
	err := lib.RegisterCRDs(ctx, certManagerCrds)
	if err != nil {
		return err
	}

	//nolint
	//// exampleCert
	err = CreateCert(ctx)
	if err != nil {
		return err
	}

	//// APP
	err = execGeneratedCode(ctx)
	if err != nil {
		return err
	}
	return nil
}

// CreateCert creates a certificate
func CreateCert(ctx *pulumi.Context) error {
	//// Instantiate a Certificate resource.
	_, err := certv1b1.NewCertificate(ctx, "example-cert", &certv1b1.CertificateArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("example-com"),
		},
		Spec: certv1b1.CertificateSpecArgs{
			SecretName:  pulumi.String("example-com-tls"),
			Duration:    pulumi.String("2160h"),
			RenewBefore: pulumi.String("360h"),
			CommonName:  pulumi.String("example.com"),
			DnsNames: pulumi.StringArray{
				pulumi.String("example.com"),
				pulumi.String("www.example.com"),
			},
			IssuerRef: certv1b1.CertificateSpecIssuerRefArgs{
				Name: pulumi.String("ca-issuer"),
				Kind: pulumi.String("Issuer"),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

//nolint
func execGeneratedCode(ctx *pulumi.Context) error {
	_, err := rbacv1.NewClusterRole(ctx, "cert_manager_cainjectorClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cainjector"),
				"app.kubernetes.io/component":  pulumi.String("cainjector"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cainjector"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-cainjector"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("create"),
					pulumi.String("update"),
					pulumi.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("admissionregistration.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("validatingwebhookconfigurations"),
					pulumi.String("mutatingwebhookconfigurations"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apiregistration.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("apiservices"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("apiextensions.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("customresourcedefinitions"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("auditregistration.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("auditsinks"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("update"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_approve_cert_manager_ioClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("cert-manager"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-approve:cert-manager-io"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("issuers.cert-manager.io/*"),
					pulumi.String("clusterissuers.cert-manager.io/*"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("signers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("approve"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_certificatesClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-certificates"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates"),
					pulumi.String("certificates/status"),
					pulumi.String("certificaterequests"),
					pulumi.String("certificaterequests/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates"),
					pulumi.String("certificaterequests"),
					pulumi.String("clusterissuers"),
					pulumi.String("issuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates/finalizers"),
					pulumi.String("certificaterequests/finalizers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("orders"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("create"),
					pulumi.String("update"),
					pulumi.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_certificatesigningrequestsClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("cert-manager"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-certificatesigningrequests"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("certificates.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificatesigningrequests"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("certificates.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificatesigningrequests/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("certificates.k8s.io"),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("issuers.cert-manager.io/*"),
					pulumi.String("clusterissuers.cert-manager.io/*"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("signers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("sign"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("authorization.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("subjectaccessreviews"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_challengesClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-challenges"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("challenges"),
					pulumi.String("challenges/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("challenges"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("issuers"),
					pulumi.String("clusterissuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("pods"),
					pulumi.String("services"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("create"),
					pulumi.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("route.openshift.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("routes/custom-host"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("challenges/finalizers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_clusterissuersClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-clusterissuers"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("clusterissuers"),
					pulumi.String("clusterissuers/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("clusterissuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("create"),
					pulumi.String("update"),
					pulumi.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_ingress_shimClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-ingress-shim"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates"),
					pulumi.String("certificaterequests"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("update"),
					pulumi.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates"),
					pulumi.String("certificaterequests"),
					pulumi.String("issuers"),
					pulumi.String("clusterissuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("networking.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("ingresses/finalizers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_issuersClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-issuers"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("issuers"),
					pulumi.String("issuers/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("issuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("create"),
					pulumi.String("update"),
					pulumi.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_ordersClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-orders"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("orders"),
					pulumi.String("orders/status"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("orders"),
					pulumi.String("challenges"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("clusterissuers"),
					pulumi.String("issuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("challenges"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("orders/finalizers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("events"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_editClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":                  pulumi.String("controller"),
				"app.kubernetes.io/instance":                   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by":                 pulumi.String("Helm"),
				"app.kubernetes.io/name":                       pulumi.String("cert-manager"),
				"helm.sh/chart":                                pulumi.String("cert-manager-v1.4.0"),
				"rbac.authorization.k8s.io/aggregate-to-admin": pulumi.String("true"),
				"rbac.authorization.k8s.io/aggregate-to-edit":  pulumi.String("true"),
			},
			Name: pulumi.String("cert-manager-edit"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates"),
					pulumi.String("certificaterequests"),
					pulumi.String("issuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("deletecollection"),
					pulumi.String("patch"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("challenges"),
					pulumi.String("orders"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
					pulumi.String("delete"),
					pulumi.String("deletecollection"),
					pulumi.String("patch"),
					pulumi.String("update"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_viewClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":                  pulumi.String("controller"),
				"app.kubernetes.io/instance":                   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by":                 pulumi.String("Helm"),
				"app.kubernetes.io/name":                       pulumi.String("cert-manager"),
				"helm.sh/chart":                                pulumi.String("cert-manager-v1.4.0"),
				"rbac.authorization.k8s.io/aggregate-to-admin": pulumi.String("true"),
				"rbac.authorization.k8s.io/aggregate-to-edit":  pulumi.String("true"),
				"rbac.authorization.k8s.io/aggregate-to-view":  pulumi.String("true"),
			},
			Name: pulumi.String("cert-manager-view"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("certificates"),
					pulumi.String("certificaterequests"),
					pulumi.String("issuers"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("acme.cert-manager.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("challenges"),
					pulumi.String("orders"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_webhook_subjectaccessreviewsClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-webhook:subjectaccessreviews"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("authorization.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("subjectaccessreviews"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_cainjectorClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cainjector"),
				"app.kubernetes.io/component":  pulumi.String("cainjector"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cainjector"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-cainjector"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-cainjector"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager-cainjector"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_approve_cert_manager_ioClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("cert-manager"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-approve:cert-manager-io"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-approve:cert-manager-io"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_certificatesClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-certificates"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-certificates"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_certificatesigningrequestsClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("cert-manager"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-certificatesigningrequests"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-certificatesigningrequests"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_challengesClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-challenges"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-challenges"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_clusterissuersClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-clusterissuers"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-clusterissuers"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_ingress_shimClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-ingress-shim"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-ingress-shim"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_issuersClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-issuers"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-issuers"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_ordersClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-controller-orders"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-controller-orders"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_webhook_subjectaccessreviewsClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-webhook:subjectaccessreviews"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("ClusterRole"),
			Name:     pulumi.String("cert-manager-webhook:subjectaccessreviews"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				ApiGroup:  pulumi.String(""),
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager-webhook"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "cert_managerCert_manager_cainjectorDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cainjector"),
				"app.kubernetes.io/component":  pulumi.String("cainjector"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cainjector"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-cainjector"),
			Namespace: pulumi.String("cert-manager"),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app.kubernetes.io/component": pulumi.String("cainjector"),
					"app.kubernetes.io/instance":  pulumi.String("cert-manager"),
					"app.kubernetes.io/name":      pulumi.String("cainjector"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app":                          pulumi.String("cainjector"),
						"app.kubernetes.io/component":  pulumi.String("cainjector"),
						"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
						"app.kubernetes.io/managed-by": pulumi.String("Helm"),
						"app.kubernetes.io/name":       pulumi.String("cainjector"),
						"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Args: pulumi.StringArray{
								pulumi.String("--v=2"),
								pulumi.String("--leader-election-namespace=kube-system"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
							},
							Image:           pulumi.String("quay.io/jetstack/cert-manager-cainjector:v1.4.0"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Name:            pulumi.String("cert-manager"),
							Resources:       nil,
						},
					},
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: pulumi.Bool(true),
					},
					ServiceAccountName: pulumi.String("cert-manager-cainjector"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "cert_managerCert_manager_webhookDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-webhook"),
			Namespace: pulumi.String("cert-manager"),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app.kubernetes.io/component": pulumi.String("webhook"),
					"app.kubernetes.io/instance":  pulumi.String("cert-manager"),
					"app.kubernetes.io/name":      pulumi.String("webhook"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app":                          pulumi.String("webhook"),
						"app.kubernetes.io/component":  pulumi.String("webhook"),
						"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
						"app.kubernetes.io/managed-by": pulumi.String("Helm"),
						"app.kubernetes.io/name":       pulumi.String("webhook"),
						"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Args: pulumi.StringArray{
								pulumi.String("--v=2"),
								pulumi.String("--secure-port=10250"),
								pulumi.String(fmt.Sprintf("%v%v%v", "--dynamic-serving-ca-secret-namespace=", "$", "(POD_NAMESPACE)")),
								pulumi.String("--dynamic-serving-ca-secret-name=cert-manager-webhook-ca"),
								pulumi.String("--dynamic-serving-dns-names=cert-manager-webhook,cert-manager-webhook.cert-manager,cert-manager-webhook.cert-manager.svc"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
							},
							Image:           pulumi.String("quay.io/jetstack/cert-manager-webhook:v1.4.0"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							LivenessProbe: &corev1.ProbeArgs{
								FailureThreshold: pulumi.Int(3),
								HttpGet: &corev1.HTTPGetActionArgs{
									Path:   pulumi.String("/livez"),
									Port:   pulumi.Int(6080),
									Scheme: pulumi.String("HTTP"),
								},
								InitialDelaySeconds: pulumi.Int(60),
								PeriodSeconds:       pulumi.Int(10),
								SuccessThreshold:    pulumi.Int(1),
								TimeoutSeconds:      pulumi.Int(1),
							},
							Name: pulumi.String("cert-manager"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(10250),
									Name:          pulumi.String("https"),
								},
							},
							ReadinessProbe: &corev1.ProbeArgs{
								FailureThreshold: pulumi.Int(3),
								HttpGet: &corev1.HTTPGetActionArgs{
									Path:   pulumi.String("/healthz"),
									Port:   pulumi.Int(6080),
									Scheme: pulumi.String("HTTP"),
								},
								InitialDelaySeconds: pulumi.Int(5),
								PeriodSeconds:       pulumi.Int(5),
								SuccessThreshold:    pulumi.Int(1),
								TimeoutSeconds:      pulumi.Int(1),
							},
							Resources: nil,
						},
					},
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: pulumi.Bool(true),
					},
					ServiceAccountName: pulumi.String("cert-manager-webhook"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "cert_managerCert_managerDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager"),
			Namespace: pulumi.String("cert-manager"),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app.kubernetes.io/component": pulumi.String("controller"),
					"app.kubernetes.io/instance":  pulumi.String("cert-manager"),
					"app.kubernetes.io/name":      pulumi.String("cert-manager"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Annotations: pulumi.StringMap{
						"prometheus.io/path":   pulumi.String("/metrics"),
						"prometheus.io/port":   pulumi.String("9402"),
						"prometheus.io/scrape": pulumi.String("true"),
					},
					Labels: pulumi.StringMap{
						"app":                          pulumi.String("cert-manager"),
						"app.kubernetes.io/component":  pulumi.String("controller"),
						"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
						"app.kubernetes.io/managed-by": pulumi.String("Helm"),
						"app.kubernetes.io/name":       pulumi.String("cert-manager"),
						"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Args: pulumi.StringArray{
								pulumi.String("--v=2"),
								pulumi.String(fmt.Sprintf("%v%v%v", "--cluster-resource-namespace=", "$", "(POD_NAMESPACE)")),
								pulumi.String("--leader-election-namespace=kube-system"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: pulumi.String("metadata.namespace"),
										},
									},
								},
							},
							Image:           pulumi.String("quay.io/jetstack/cert-manager-controller:v1.4.0"),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Name:            pulumi.String("cert-manager"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(9402),
									Protocol:      pulumi.String("TCP"),
								},
							},
							Resources: nil,
						},
					},
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: pulumi.Bool(true),
					},
					ServiceAccountName: pulumi.String("cert-manager"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = admissionregistrationv1.NewMutatingWebhookConfiguration(ctx, "cert_manager_webhookMutatingWebhookConfiguration", &admissionregistrationv1.MutatingWebhookConfigurationArgs{
		ApiVersion: pulumi.String("admissionregistration.k8s.io/v1"),
		Kind:       pulumi.String("MutatingWebhookConfiguration"),
		Metadata: &metav1.ObjectMetaArgs{
			Annotations: pulumi.StringMap{
				"cert-manager.io/inject-ca-from-secret": pulumi.String("cert-manager/cert-manager-webhook-ca"),
			},
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-webhook"),
		},
		Webhooks: admissionregistrationv1.MutatingWebhookArray{
			&admissionregistrationv1.MutatingWebhookArgs{
				AdmissionReviewVersions: pulumi.StringArray{
					pulumi.String("v1"),
					pulumi.String("v1beta1"),
				},
				ClientConfig: &admissionregistrationv1.WebhookClientConfigArgs{
					Service: &admissionregistrationv1.ServiceReferenceArgs{
						Name:      pulumi.String("cert-manager-webhook"),
						Namespace: pulumi.String("cert-manager"),
						Path:      pulumi.String("/mutate"),
					},
				},
				FailurePolicy: pulumi.String("Fail"),
				Name:          pulumi.String("webhook.cert-manager.io"),
				Rules: admissionregistrationv1.RuleWithOperationsArray{
					&admissionregistrationv1.RuleWithOperationsArgs{
						ApiGroups: pulumi.StringArray{
							pulumi.String("cert-manager.io"),
							pulumi.String("acme.cert-manager.io"),
						},
						ApiVersions: pulumi.StringArray{
							pulumi.String("*"),
						},
						Operations: pulumi.StringArray{
							pulumi.String("CREATE"),
							pulumi.String("UPDATE"),
						},
						Resources: pulumi.StringArray{
							pulumi.String("*/*"),
						},
					},
				},
				SideEffects:    pulumi.String("None"),
				TimeoutSeconds: pulumi.Int(10),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewNamespace(ctx, "cert_managerNamespace", &corev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Namespace"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRole(ctx, "kube_systemCert_manager_cainjector_leaderelectionRole", &rbacv1.RoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cainjector"),
				"app.kubernetes.io/component":  pulumi.String("cainjector"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cainjector"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-cainjector:leaderelection"),
			Namespace: pulumi.String("kube-system"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("cert-manager-cainjector-leader-election"),
					pulumi.String("cert-manager-cainjector-leader-election-core"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("update"),
					pulumi.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("coordination.k8s.io"),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("cert-manager-cainjector-leader-election"),
					pulumi.String("cert-manager-cainjector-leader-election-core"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("leases"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("update"),
					pulumi.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("coordination.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("leases"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRole(ctx, "cert_managerCert_manager_webhook_dynamic_servingRole", &rbacv1.RoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-webhook:dynamic-serving"),
			Namespace: pulumi.String("cert-manager"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("cert-manager-webhook-ca"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("list"),
					pulumi.String("watch"),
					pulumi.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("secrets"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRole(ctx, "kube_systemCert_manager_leaderelectionRole", &rbacv1.RoleArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager:leaderelection"),
			Namespace: pulumi.String("kube-system"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("cert-manager-controller"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("update"),
					pulumi.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String(""),
				},
				Resources: pulumi.StringArray{
					pulumi.String("configmaps"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("coordination.k8s.io"),
				},
				ResourceNames: pulumi.StringArray{
					pulumi.String("cert-manager-controller"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("leases"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("get"),
					pulumi.String("update"),
					pulumi.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: pulumi.StringArray{
					pulumi.String("coordination.k8s.io"),
				},
				Resources: pulumi.StringArray{
					pulumi.String("leases"),
				},
				Verbs: pulumi.StringArray{
					pulumi.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "kube_systemCert_manager_cainjector_leaderelectionRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cainjector"),
				"app.kubernetes.io/component":  pulumi.String("cainjector"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cainjector"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-cainjector:leaderelection"),
			Namespace: pulumi.String("kube-system"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("Role"),
			Name:     pulumi.String("cert-manager-cainjector:leaderelection"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager-cainjector"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "cert_managerCert_manager_webhook_dynamic_servingRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-webhook:dynamic-serving"),
			Namespace: pulumi.String("cert-manager"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("Role"),
			Name:     pulumi.String("cert-manager-webhook:dynamic-serving"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				ApiGroup:  pulumi.String(""),
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager-webhook"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "kube_systemCert_manager_leaderelectionRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: pulumi.String("rbac.authorization.k8s.io/v1"),
		Kind:       pulumi.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager:leaderelection"),
			Namespace: pulumi.String("kube-system"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
			Kind:     pulumi.String("Role"),
			Name:     pulumi.String("cert-manager:leaderelection"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				ApiGroup:  pulumi.String(""),
				Kind:      pulumi.String("ServiceAccount"),
				Name:      pulumi.String("cert-manager"),
				Namespace: pulumi.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "cert_managerCert_manager_webhookService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-webhook"),
			Namespace: pulumi.String("cert-manager"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("https"),
					Port:       pulumi.Int(443),
					TargetPort: pulumi.Int(10250),
				},
			},
			Selector: pulumi.StringMap{
				"app.kubernetes.io/component": pulumi.String("webhook"),
				"app.kubernetes.io/instance":  pulumi.String("cert-manager"),
				"app.kubernetes.io/name":      pulumi.String("webhook"),
			},
			Type: pulumi.String("ClusterIP"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "cert_managerCert_managerService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager"),
			Namespace: pulumi.String("cert-manager"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(9402),
					Protocol:   pulumi.String("TCP"),
					TargetPort: pulumi.Int(9402),
				},
			},
			Selector: pulumi.StringMap{
				"app.kubernetes.io/component": pulumi.String("controller"),
				"app.kubernetes.io/instance":  pulumi.String("cert-manager"),
				"app.kubernetes.io/name":      pulumi.String("cert-manager"),
			},
			Type: pulumi.String("ClusterIP"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "cert_managerCert_manager_cainjectorServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion:                   pulumi.String("v1"),
		AutomountServiceAccountToken: pulumi.Bool(true),
		Kind:                         pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cainjector"),
				"app.kubernetes.io/component":  pulumi.String("cainjector"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cainjector"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-cainjector"),
			Namespace: pulumi.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "cert_managerCert_manager_webhookServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion:                   pulumi.String("v1"),
		AutomountServiceAccountToken: pulumi.Bool(true),
		Kind:                         pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager-webhook"),
			Namespace: pulumi.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "cert_managerCert_managerServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion:                   pulumi.String("v1"),
		AutomountServiceAccountToken: pulumi.Bool(true),
		Kind:                         pulumi.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("cert-manager"),
				"app.kubernetes.io/component":  pulumi.String("controller"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("cert-manager"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name:      pulumi.String("cert-manager"),
			Namespace: pulumi.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = admissionregistrationv1.NewValidatingWebhookConfiguration(ctx, "cert_manager_webhookValidatingWebhookConfiguration", &admissionregistrationv1.ValidatingWebhookConfigurationArgs{
		ApiVersion: pulumi.String("admissionregistration.k8s.io/v1"),
		Kind:       pulumi.String("ValidatingWebhookConfiguration"),
		Metadata: &metav1.ObjectMetaArgs{
			Annotations: pulumi.StringMap{
				"cert-manager.io/inject-ca-from-secret": pulumi.String("cert-manager/cert-manager-webhook-ca"),
			},
			Labels: pulumi.StringMap{
				"app":                          pulumi.String("webhook"),
				"app.kubernetes.io/component":  pulumi.String("webhook"),
				"app.kubernetes.io/instance":   pulumi.String("cert-manager"),
				"app.kubernetes.io/managed-by": pulumi.String("Helm"),
				"app.kubernetes.io/name":       pulumi.String("webhook"),
				"helm.sh/chart":                pulumi.String("cert-manager-v1.4.0"),
			},
			Name: pulumi.String("cert-manager-webhook"),
		},
		Webhooks: admissionregistrationv1.ValidatingWebhookArray{
			&admissionregistrationv1.ValidatingWebhookArgs{
				AdmissionReviewVersions: pulumi.StringArray{
					pulumi.String("v1"),
					pulumi.String("v1beta1"),
				},
				ClientConfig: &admissionregistrationv1.WebhookClientConfigArgs{
					Service: &admissionregistrationv1.ServiceReferenceArgs{
						Name:      pulumi.String("cert-manager-webhook"),
						Namespace: pulumi.String("cert-manager"),
						Path:      pulumi.String("/validate"),
					},
				},
				FailurePolicy: pulumi.String("Fail"),
				Name:          pulumi.String("webhook.cert-manager.io"),
				NamespaceSelector: &metav1.LabelSelectorArgs{
					MatchExpressions: metav1.LabelSelectorRequirementArray{
						&metav1.LabelSelectorRequirementArgs{
							Key:      pulumi.String("cert-manager.io/disable-validation"),
							Operator: pulumi.String("NotIn"),
							Values: pulumi.StringArray{
								pulumi.String("true"),
							},
						},
						&metav1.LabelSelectorRequirementArgs{
							Key:      pulumi.String("name"),
							Operator: pulumi.String("NotIn"),
							Values: pulumi.StringArray{
								pulumi.String("cert-manager"),
							},
						},
					},
				},
				Rules: admissionregistrationv1.RuleWithOperationsArray{
					&admissionregistrationv1.RuleWithOperationsArgs{
						ApiGroups: pulumi.StringArray{
							pulumi.String("cert-manager.io"),
							pulumi.String("acme.cert-manager.io"),
						},
						ApiVersions: pulumi.StringArray{
							pulumi.String("*"),
						},
						Operations: pulumi.StringArray{
							pulumi.String("CREATE"),
							pulumi.String("UPDATE"),
						},
						Resources: pulumi.StringArray{
							pulumi.String("*/*"),
						},
					},
				},
				SideEffects:    pulumi.String("None"),
				TimeoutSeconds: pulumi.Int(10),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
