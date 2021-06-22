package certmanager

import (
	"fmt"

	admissionregistrationv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/admissionregistration/v1"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/rbac/v1"
	p "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	certv1b1 "thesym.site/k8s/crds/cert-manager/certmanager/v1beta1"
)

func CreateCertmanager(ctx *p.Context) error {

	certManagerCrds := []Crd{
		{"certmanager-order-definition", "./crds/cert-manager/cdrDefinitions/customresourcedefinition-orders.acme.cert-manager.io.yaml"},
		{"certmanager-issuer-definition", "./crds/cert-manager/cdrDefinitions/customresourcedefinition-issuers.cert-manager.io.yaml"},
		{"certmanager-challenge-definition", "./crds/cert-manager/cdrDefinitions/customresourcedefinition-challenges.acme.cert-manager.io.yaml"},
		{"certmanager-certificate-definition", "./crds/cert-manager/cdrDefinitions/customresourcedefinition-certificates.cert-manager.io.yaml"},
		{"certmanager-clusterIssuer-definition", "./crds/cert-manager/cdrDefinitions/customresourcedefinition-clusterissuers.cert-manager.io.yaml"},
		{"certmanager-certificate-request-definition", "./crds/cert-manager/cdrDefinitions/customresourcedefinition-certificaterequests.cert-manager.io.yaml"},
	}
	//// Register the CRD.
	err := registerCRDs(ctx, certManagerCrds)
	if err != nil {
		return err
	}

	//// Instantiate a Certificate resource.
	_, err = certv1b1.NewCertificate(ctx, "example-cert", &certv1b1.CertificateArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("example-com"),
		},
		Spec: certv1b1.CertificateSpecArgs{
			SecretName:  p.String("example-com-tls"),
			Duration:    p.String("2160h"),
			RenewBefore: p.String("360h"),
			CommonName:  p.String("example.com"),
			DnsNames: p.StringArray{
				p.String("example.com"),
				p.String("www.example.com"),
			},
			IssuerRef: certv1b1.CertificateSpecIssuerRefArgs{
				Name: p.String("ca-issuer"),
				Kind: p.String("Issuer"),
			},
		},
	})
	if err != nil {
		return err
	}

	//// APP
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_cainjectorClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cainjector"),
				"app.kubernetes.io/component":  p.String("cainjector"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cainjector"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-cainjector"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("events"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("create"),
					p.String("update"),
					p.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("admissionregistration.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("validatingwebhookconfigurations"),
					p.String("mutatingwebhookconfigurations"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("apiregistration.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("apiservices"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("apiextensions.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("customresourcedefinitions"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("auditregistration.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("auditsinks"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("update"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_approve_cert_manager_ioClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("cert-manager"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-approve:cert-manager-io"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				ResourceNames: p.StringArray{
					p.String("issuers.cert-manager.io/*"),
					p.String("clusterissuers.cert-manager.io/*"),
				},
				Resources: p.StringArray{
					p.String("signers"),
				},
				Verbs: p.StringArray{
					p.String("approve"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_certificatesClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-certificates"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates"),
					p.String("certificates/status"),
					p.String("certificaterequests"),
					p.String("certificaterequests/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates"),
					p.String("certificaterequests"),
					p.String("clusterissuers"),
					p.String("issuers"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates/finalizers"),
					p.String("certificaterequests/finalizers"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("orders"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("delete"),
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("create"),
					p.String("update"),
					p.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("events"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_certificatesigningrequestsClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("cert-manager"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-certificatesigningrequests"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("certificates.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("certificatesigningrequests"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("certificates.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("certificatesigningrequests/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("certificates.k8s.io"),
				},
				ResourceNames: p.StringArray{
					p.String("issuers.cert-manager.io/*"),
					p.String("clusterissuers.cert-manager.io/*"),
				},
				Resources: p.StringArray{
					p.String("signers"),
				},
				Verbs: p.StringArray{
					p.String("sign"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("authorization.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("subjectaccessreviews"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_challengesClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-challenges"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("challenges"),
					p.String("challenges/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("challenges"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("issuers"),
					p.String("clusterissuers"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("events"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("pods"),
					p.String("services"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("create"),
					p.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("networking.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("ingresses"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("create"),
					p.String("delete"),
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("route.openshift.io"),
				},
				Resources: p.StringArray{
					p.String("routes/custom-host"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("challenges/finalizers"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_clusterissuersClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-clusterissuers"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("clusterissuers"),
					p.String("clusterissuers/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("clusterissuers"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("create"),
					p.String("update"),
					p.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("events"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_ingress_shimClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-ingress-shim"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates"),
					p.String("certificaterequests"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("update"),
					p.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates"),
					p.String("certificaterequests"),
					p.String("issuers"),
					p.String("clusterissuers"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("networking.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("ingresses"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("networking.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("ingresses/finalizers"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("events"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_issuersClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-issuers"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("issuers"),
					p.String("issuers/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("issuers"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("create"),
					p.String("update"),
					p.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("events"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_controller_ordersClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-orders"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("orders"),
					p.String("orders/status"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("orders"),
					p.String("challenges"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("clusterissuers"),
					p.String("issuers"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("challenges"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("delete"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("orders/finalizers"),
				},
				Verbs: p.StringArray{
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("events"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("patch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_editClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                                          p.String("cert-manager"),
				"app.kubernetes.io/component":                  p.String("controller"),
				"app.kubernetes.io/instance":                   p.String("cert-manager"),
				"app.kubernetes.io/managed-by":                 p.String("Helm"),
				"app.kubernetes.io/name":                       p.String("cert-manager"),
				"helm.sh/chart":                                p.String("cert-manager-v1.4.0"),
				"rbac.authorization.k8s.io/aggregate-to-admin": p.String("true"),
				"rbac.authorization.k8s.io/aggregate-to-edit":  p.String("true"),
			},
			Name: p.String("cert-manager-edit"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates"),
					p.String("certificaterequests"),
					p.String("issuers"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("delete"),
					p.String("deletecollection"),
					p.String("patch"),
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("challenges"),
					p.String("orders"),
				},
				Verbs: p.StringArray{
					p.String("create"),
					p.String("delete"),
					p.String("deletecollection"),
					p.String("patch"),
					p.String("update"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_viewClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                                          p.String("cert-manager"),
				"app.kubernetes.io/component":                  p.String("controller"),
				"app.kubernetes.io/instance":                   p.String("cert-manager"),
				"app.kubernetes.io/managed-by":                 p.String("Helm"),
				"app.kubernetes.io/name":                       p.String("cert-manager"),
				"helm.sh/chart":                                p.String("cert-manager-v1.4.0"),
				"rbac.authorization.k8s.io/aggregate-to-admin": p.String("true"),
				"rbac.authorization.k8s.io/aggregate-to-edit":  p.String("true"),
				"rbac.authorization.k8s.io/aggregate-to-view":  p.String("true"),
			},
			Name: p.String("cert-manager-view"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("certificates"),
					p.String("certificaterequests"),
					p.String("issuers"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("acme.cert-manager.io"),
				},
				Resources: p.StringArray{
					p.String("challenges"),
					p.String("orders"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRole(ctx, "cert_manager_webhook_subjectaccessreviewsClusterRole", &rbacv1.ClusterRoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRole"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-webhook:subjectaccessreviews"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("authorization.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("subjectaccessreviews"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_cainjectorClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cainjector"),
				"app.kubernetes.io/component":  p.String("cainjector"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cainjector"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-cainjector"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-cainjector"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager-cainjector"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_approve_cert_manager_ioClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("cert-manager"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-approve:cert-manager-io"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-approve:cert-manager-io"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_certificatesClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-certificates"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-certificates"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_certificatesigningrequestsClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("cert-manager"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-certificatesigningrequests"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-certificatesigningrequests"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_challengesClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-challenges"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-challenges"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_clusterissuersClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-clusterissuers"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-clusterissuers"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_ingress_shimClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-ingress-shim"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-ingress-shim"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_issuersClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-issuers"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-issuers"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_controller_ordersClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-controller-orders"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-controller-orders"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewClusterRoleBinding(ctx, "cert_manager_webhook_subjectaccessreviewsClusterRoleBinding", &rbacv1.ClusterRoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("ClusterRoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-webhook:subjectaccessreviews"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("ClusterRole"),
			Name:     p.String("cert-manager-webhook:subjectaccessreviews"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				ApiGroup:  p.String(""),
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager-webhook"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "cert_managerCert_manager_cainjectorDeployment", &appsv1.DeploymentArgs{
		ApiVersion: p.String("apps/v1"),
		Kind:       p.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cainjector"),
				"app.kubernetes.io/component":  p.String("cainjector"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cainjector"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-cainjector"),
			Namespace: p.String("cert-manager"),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: p.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: p.StringMap{
					"app.kubernetes.io/component": p.String("cainjector"),
					"app.kubernetes.io/instance":  p.String("cert-manager"),
					"app.kubernetes.io/name":      p.String("cainjector"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: p.StringMap{
						"app":                          p.String("cainjector"),
						"app.kubernetes.io/component":  p.String("cainjector"),
						"app.kubernetes.io/instance":   p.String("cert-manager"),
						"app.kubernetes.io/managed-by": p.String("Helm"),
						"app.kubernetes.io/name":       p.String("cainjector"),
						"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Args: p.StringArray{
								p.String("--v=2"),
								p.String("--leader-election-namespace=kube-system"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: p.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: p.String("metadata.namespace"),
										},
									},
								},
							},
							Image:           p.String("quay.io/jetstack/cert-manager-cainjector:v1.4.0"),
							ImagePullPolicy: p.String("IfNotPresent"),
							Name:            p.String("cert-manager"),
							Resources:       nil,
						},
					},
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: p.Bool(true),
					},
					ServiceAccountName: p.String("cert-manager-cainjector"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "cert_managerCert_manager_webhookDeployment", &appsv1.DeploymentArgs{
		ApiVersion: p.String("apps/v1"),
		Kind:       p.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-webhook"),
			Namespace: p.String("cert-manager"),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: p.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: p.StringMap{
					"app.kubernetes.io/component": p.String("webhook"),
					"app.kubernetes.io/instance":  p.String("cert-manager"),
					"app.kubernetes.io/name":      p.String("webhook"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: p.StringMap{
						"app":                          p.String("webhook"),
						"app.kubernetes.io/component":  p.String("webhook"),
						"app.kubernetes.io/instance":   p.String("cert-manager"),
						"app.kubernetes.io/managed-by": p.String("Helm"),
						"app.kubernetes.io/name":       p.String("webhook"),
						"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Args: p.StringArray{
								p.String("--v=2"),
								p.String("--secure-port=10250"),
								p.String(fmt.Sprintf("%v%v%v", "--dynamic-serving-ca-secret-namespace=", "$", "(POD_NAMESPACE)")),
								p.String("--dynamic-serving-ca-secret-name=cert-manager-webhook-ca"),
								p.String("--dynamic-serving-dns-names=cert-manager-webhook,cert-manager-webhook.cert-manager,cert-manager-webhook.cert-manager.svc"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: p.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: p.String("metadata.namespace"),
										},
									},
								},
							},
							Image:           p.String("quay.io/jetstack/cert-manager-webhook:v1.4.0"),
							ImagePullPolicy: p.String("IfNotPresent"),
							LivenessProbe: &corev1.ProbeArgs{
								FailureThreshold: p.Int(3),
								HttpGet: &corev1.HTTPGetActionArgs{
									Path:   p.String("/livez"),
									Port:   p.Int(6080),
									Scheme: p.String("HTTP"),
								},
								InitialDelaySeconds: p.Int(60),
								PeriodSeconds:       p.Int(10),
								SuccessThreshold:    p.Int(1),
								TimeoutSeconds:      p.Int(1),
							},
							Name: p.String("cert-manager"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: p.Int(10250),
									Name:          p.String("https"),
								},
							},
							ReadinessProbe: &corev1.ProbeArgs{
								FailureThreshold: p.Int(3),
								HttpGet: &corev1.HTTPGetActionArgs{
									Path:   p.String("/healthz"),
									Port:   p.Int(6080),
									Scheme: p.String("HTTP"),
								},
								InitialDelaySeconds: p.Int(5),
								PeriodSeconds:       p.Int(5),
								SuccessThreshold:    p.Int(1),
								TimeoutSeconds:      p.Int(1),
							},
							Resources: nil,
						},
					},
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: p.Bool(true),
					},
					ServiceAccountName: p.String("cert-manager-webhook"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "cert_managerCert_managerDeployment", &appsv1.DeploymentArgs{
		ApiVersion: p.String("apps/v1"),
		Kind:       p.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager"),
			Namespace: p.String("cert-manager"),
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: p.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: p.StringMap{
					"app.kubernetes.io/component": p.String("controller"),
					"app.kubernetes.io/instance":  p.String("cert-manager"),
					"app.kubernetes.io/name":      p.String("cert-manager"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Annotations: p.StringMap{
						"prometheus.io/path":   p.String("/metrics"),
						"prometheus.io/port":   p.String("9402"),
						"prometheus.io/scrape": p.String("true"),
					},
					Labels: p.StringMap{
						"app":                          p.String("cert-manager"),
						"app.kubernetes.io/component":  p.String("controller"),
						"app.kubernetes.io/instance":   p.String("cert-manager"),
						"app.kubernetes.io/managed-by": p.String("Helm"),
						"app.kubernetes.io/name":       p.String("cert-manager"),
						"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Args: p.StringArray{
								p.String("--v=2"),
								p.String(fmt.Sprintf("%v%v%v", "--cluster-resource-namespace=", "$", "(POD_NAMESPACE)")),
								p.String("--leader-election-namespace=kube-system"),
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: p.String("POD_NAMESPACE"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										FieldRef: &corev1.ObjectFieldSelectorArgs{
											FieldPath: p.String("metadata.namespace"),
										},
									},
								},
							},
							Image:           p.String("quay.io/jetstack/cert-manager-controller:v1.4.0"),
							ImagePullPolicy: p.String("IfNotPresent"),
							Name:            p.String("cert-manager"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: p.Int(9402),
									Protocol:      p.String("TCP"),
								},
							},
							Resources: nil,
						},
					},
					SecurityContext: &corev1.PodSecurityContextArgs{
						RunAsNonRoot: p.Bool(true),
					},
					ServiceAccountName: p.String("cert-manager"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = admissionregistrationv1.NewMutatingWebhookConfiguration(ctx, "cert_manager_webhookMutatingWebhookConfiguration", &admissionregistrationv1.MutatingWebhookConfigurationArgs{
		ApiVersion: p.String("admissionregistration.k8s.io/v1"),
		Kind:       p.String("MutatingWebhookConfiguration"),
		Metadata: &metav1.ObjectMetaArgs{
			Annotations: p.StringMap{
				"cert-manager.io/inject-ca-from-secret": p.String("cert-manager/cert-manager-webhook-ca"),
			},
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-webhook"),
		},
		Webhooks: admissionregistrationv1.MutatingWebhookArray{
			&admissionregistrationv1.MutatingWebhookArgs{
				AdmissionReviewVersions: p.StringArray{
					p.String("v1"),
					p.String("v1beta1"),
				},
				ClientConfig: &admissionregistrationv1.WebhookClientConfigArgs{
					Service: &admissionregistrationv1.ServiceReferenceArgs{
						Name:      p.String("cert-manager-webhook"),
						Namespace: p.String("cert-manager"),
						Path:      p.String("/mutate"),
					},
				},
				FailurePolicy: p.String("Fail"),
				Name:          p.String("webhook.cert-manager.io"),
				Rules: admissionregistrationv1.RuleWithOperationsArray{
					&admissionregistrationv1.RuleWithOperationsArgs{
						ApiGroups: p.StringArray{
							p.String("cert-manager.io"),
							p.String("acme.cert-manager.io"),
						},
						ApiVersions: p.StringArray{
							p.String("*"),
						},
						Operations: p.StringArray{
							p.String("CREATE"),
							p.String("UPDATE"),
						},
						Resources: p.StringArray{
							p.String("*/*"),
						},
					},
				},
				SideEffects:    p.String("None"),
				TimeoutSeconds: p.Int(10),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewNamespace(ctx, "cert_managerNamespace", &corev1.NamespaceArgs{
		ApiVersion: p.String("v1"),
		Kind:       p.String("Namespace"),
		Metadata: &metav1.ObjectMetaArgs{
			Name: p.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRole(ctx, "kube_systemCert_manager_cainjector_leaderelectionRole", &rbacv1.RoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cainjector"),
				"app.kubernetes.io/component":  p.String("cainjector"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cainjector"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-cainjector:leaderelection"),
			Namespace: p.String("kube-system"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				ResourceNames: p.StringArray{
					p.String("cert-manager-cainjector-leader-election"),
					p.String("cert-manager-cainjector-leader-election-core"),
				},
				Resources: p.StringArray{
					p.String("configmaps"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("update"),
					p.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("configmaps"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("coordination.k8s.io"),
				},
				ResourceNames: p.StringArray{
					p.String("cert-manager-cainjector-leader-election"),
					p.String("cert-manager-cainjector-leader-election-core"),
				},
				Resources: p.StringArray{
					p.String("leases"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("update"),
					p.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("coordination.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("leases"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRole(ctx, "cert_managerCert_manager_webhook_dynamic_servingRole", &rbacv1.RoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-webhook:dynamic-serving"),
			Namespace: p.String("cert-manager"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				ResourceNames: p.StringArray{
					p.String("cert-manager-webhook-ca"),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("list"),
					p.String("watch"),
					p.String("update"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("secrets"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRole(ctx, "kube_systemCert_manager_leaderelectionRole", &rbacv1.RoleArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("Role"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager:leaderelection"),
			Namespace: p.String("kube-system"),
		},
		Rules: rbacv1.PolicyRuleArray{
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				ResourceNames: p.StringArray{
					p.String("cert-manager-controller"),
				},
				Resources: p.StringArray{
					p.String("configmaps"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("update"),
					p.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String(""),
				},
				Resources: p.StringArray{
					p.String("configmaps"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("coordination.k8s.io"),
				},
				ResourceNames: p.StringArray{
					p.String("cert-manager-controller"),
				},
				Resources: p.StringArray{
					p.String("leases"),
				},
				Verbs: p.StringArray{
					p.String("get"),
					p.String("update"),
					p.String("patch"),
				},
			},
			&rbacv1.PolicyRuleArgs{
				ApiGroups: p.StringArray{
					p.String("coordination.k8s.io"),
				},
				Resources: p.StringArray{
					p.String("leases"),
				},
				Verbs: p.StringArray{
					p.String("create"),
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "kube_systemCert_manager_cainjector_leaderelectionRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cainjector"),
				"app.kubernetes.io/component":  p.String("cainjector"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cainjector"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-cainjector:leaderelection"),
			Namespace: p.String("kube-system"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("Role"),
			Name:     p.String("cert-manager-cainjector:leaderelection"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager-cainjector"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "cert_managerCert_manager_webhook_dynamic_servingRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-webhook:dynamic-serving"),
			Namespace: p.String("cert-manager"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("Role"),
			Name:     p.String("cert-manager-webhook:dynamic-serving"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				ApiGroup:  p.String(""),
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager-webhook"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = rbacv1.NewRoleBinding(ctx, "kube_systemCert_manager_leaderelectionRoleBinding", &rbacv1.RoleBindingArgs{
		ApiVersion: p.String("rbac.authorization.k8s.io/v1"),
		Kind:       p.String("RoleBinding"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager:leaderelection"),
			Namespace: p.String("kube-system"),
		},
		RoleRef: &rbacv1.RoleRefArgs{
			ApiGroup: p.String("rbac.authorization.k8s.io"),
			Kind:     p.String("Role"),
			Name:     p.String("cert-manager:leaderelection"),
		},
		Subjects: rbacv1.SubjectArray{
			&rbacv1.SubjectArgs{
				ApiGroup:  p.String(""),
				Kind:      p.String("ServiceAccount"),
				Name:      p.String("cert-manager"),
				Namespace: p.String("cert-manager"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "cert_managerCert_manager_webhookService", &corev1.ServiceArgs{
		ApiVersion: p.String("v1"),
		Kind:       p.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-webhook"),
			Namespace: p.String("cert-manager"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       p.String("https"),
					Port:       p.Int(443),
					TargetPort: p.Int(10250),
				},
			},
			Selector: p.StringMap{
				"app.kubernetes.io/component": p.String("webhook"),
				"app.kubernetes.io/instance":  p.String("cert-manager"),
				"app.kubernetes.io/name":      p.String("webhook"),
			},
			Type: p.String("ClusterIP"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "cert_managerCert_managerService", &corev1.ServiceArgs{
		ApiVersion: p.String("v1"),
		Kind:       p.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager"),
			Namespace: p.String("cert-manager"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       p.Int(9402),
					Protocol:   p.String("TCP"),
					TargetPort: p.Int(9402),
				},
			},
			Selector: p.StringMap{
				"app.kubernetes.io/component": p.String("controller"),
				"app.kubernetes.io/instance":  p.String("cert-manager"),
				"app.kubernetes.io/name":      p.String("cert-manager"),
			},
			Type: p.String("ClusterIP"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "cert_managerCert_manager_cainjectorServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion:                   p.String("v1"),
		AutomountServiceAccountToken: p.Bool(true),
		Kind:                         p.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cainjector"),
				"app.kubernetes.io/component":  p.String("cainjector"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cainjector"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-cainjector"),
			Namespace: p.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "cert_managerCert_manager_webhookServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion:                   p.String("v1"),
		AutomountServiceAccountToken: p.Bool(true),
		Kind:                         p.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager-webhook"),
			Namespace: p.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewServiceAccount(ctx, "cert_managerCert_managerServiceAccount", &corev1.ServiceAccountArgs{
		ApiVersion:                   p.String("v1"),
		AutomountServiceAccountToken: p.Bool(true),
		Kind:                         p.String("ServiceAccount"),
		Metadata: &metav1.ObjectMetaArgs{
			Labels: p.StringMap{
				"app":                          p.String("cert-manager"),
				"app.kubernetes.io/component":  p.String("controller"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("cert-manager"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name:      p.String("cert-manager"),
			Namespace: p.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}
	_, err = admissionregistrationv1.NewValidatingWebhookConfiguration(ctx, "cert_manager_webhookValidatingWebhookConfiguration", &admissionregistrationv1.ValidatingWebhookConfigurationArgs{
		ApiVersion: p.String("admissionregistration.k8s.io/v1"),
		Kind:       p.String("ValidatingWebhookConfiguration"),
		Metadata: &metav1.ObjectMetaArgs{
			Annotations: p.StringMap{
				"cert-manager.io/inject-ca-from-secret": p.String("cert-manager/cert-manager-webhook-ca"),
			},
			Labels: p.StringMap{
				"app":                          p.String("webhook"),
				"app.kubernetes.io/component":  p.String("webhook"),
				"app.kubernetes.io/instance":   p.String("cert-manager"),
				"app.kubernetes.io/managed-by": p.String("Helm"),
				"app.kubernetes.io/name":       p.String("webhook"),
				"helm.sh/chart":                p.String("cert-manager-v1.4.0"),
			},
			Name: p.String("cert-manager-webhook"),
		},
		Webhooks: admissionregistrationv1.ValidatingWebhookArray{
			&admissionregistrationv1.ValidatingWebhookArgs{
				AdmissionReviewVersions: p.StringArray{
					p.String("v1"),
					p.String("v1beta1"),
				},
				ClientConfig: &admissionregistrationv1.WebhookClientConfigArgs{
					Service: &admissionregistrationv1.ServiceReferenceArgs{
						Name:      p.String("cert-manager-webhook"),
						Namespace: p.String("cert-manager"),
						Path:      p.String("/validate"),
					},
				},
				FailurePolicy: p.String("Fail"),
				Name:          p.String("webhook.cert-manager.io"),
				NamespaceSelector: &metav1.LabelSelectorArgs{
					MatchExpressions: metav1.LabelSelectorRequirementArray{
						&metav1.LabelSelectorRequirementArgs{
							Key:      p.String("cert-manager.io/disable-validation"),
							Operator: p.String("NotIn"),
							Values: p.StringArray{
								p.String("true"),
							},
						},
						&metav1.LabelSelectorRequirementArgs{
							Key:      p.String("name"),
							Operator: p.String("NotIn"),
							Values: p.StringArray{
								p.String("cert-manager"),
							},
						},
					},
				},
				Rules: admissionregistrationv1.RuleWithOperationsArray{
					&admissionregistrationv1.RuleWithOperationsArgs{
						ApiGroups: p.StringArray{
							p.String("cert-manager.io"),
							p.String("acme.cert-manager.io"),
						},
						ApiVersions: p.StringArray{
							p.String("*"),
						},
						Operations: p.StringArray{
							p.String("CREATE"),
							p.String("UPDATE"),
						},
						Resources: p.StringArray{
							p.String("*/*"),
						},
					},
				},
				SideEffects:    p.String("None"),
				TimeoutSeconds: p.Int(10),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
