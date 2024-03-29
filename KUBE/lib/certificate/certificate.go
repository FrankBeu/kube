package certificate

import (
	"errors"
	"strconv"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	certv1 "thesym.site/kube/crds/kubernetes/certmanager/v1"
	"thesym.site/kube/lib/kubeconfig"
	"thesym.site/kube/lib/types"
)

func CreateCert(ctx *pulumi.Context, cert *types.Cert) (*certv1.Certificate, error) {
	if cert.Duration == "" {
		cert.Duration = strconv.Itoa(types.CertificateDefaultDurationInDays*24) + "h" //nolint:gomnd
	}

	domainNameSuffix := kubeconfig.DomainNameSuffix(ctx)

	commonName := ""
	if cert.CommonNameSegment == "" {
		commonName = cert.Name + domainNameSuffix
	} else {
		commonName = cert.CommonNameSegment + domainNameSuffix
	}

	dnsNames := pulumi.StringArray{
		pulumi.String(commonName),
	}
	for _, subdomainName := range cert.AdditionalSubdomainNames {
		dnsNames = append(dnsNames, pulumi.String(subdomainName+domainNameSuffix))
	}

	certificate, err := certv1.NewCertificate(ctx, cert.Name, &certv1.CertificateArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(cert.Name),
			Namespace: pulumi.String(cert.Namespace),
		},
		Spec: certv1.CertificateSpecArgs{
			SecretName:  pulumi.String(cert.Name + "-tls"),
			Duration:    pulumi.String(cert.Duration),
			RenewBefore: pulumi.String("360h"),
			CommonName:  pulumi.String(commonName),
			DnsNames:    dnsNames,
			IssuerRef: certv1.CertificateSpecIssuerRefArgs{
				Name: pulumi.String(cert.ClusterIssuerType.String()),
				//// currently only ClusterIssuers are used
				Kind: pulumi.String("ClusterIssuer"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

// CreateClusterIssuer creates an internal clusterIssuer with a localCA~ or an letsencrypt~based issuer.
func CreateClusterIssuer(ctx *pulumi.Context, clusterIssuerType types.ClusterIssuerType) (*certv1.ClusterIssuer, error) {
	adminEmail := kubeconfig.AdminEmail(ctx)

	if clusterIssuerType == types.ClusterIssuerTypeCALocal {
		err := createCALocalSecret(ctx, clusterIssuerType)
		if err != nil {
			return nil, err
		}
	}

	clusterIssuer, err := createClusterIssuer(ctx, clusterIssuerType, adminEmail)
	if err != nil {
		return nil, err
	}
	return clusterIssuer, nil
}

func createClusterIssuer(ctx *pulumi.Context, clusterIssuerType types.ClusterIssuerType, adminEmail string) (*certv1.ClusterIssuer, error) {
	var clusterIssuerSpec *certv1.ClusterIssuerSpecArgs

	if clusterIssuerType == types.ClusterIssuerTypeCALocal {
		clusterIssuerSpec = &certv1.ClusterIssuerSpecArgs{
			Ca: &certv1.ClusterIssuerSpecCaArgs{
				// SecretName: pulumi.String(clusterIssuerSecretName),
				SecretName: pulumi.String(clusterIssuerType.String()),
			},
		}
	} else {
		acmeServerURL, err := acmeServerURL(clusterIssuerType)
		if err != nil {
			return nil, err
		}

		clusterIssuerSpec = &certv1.ClusterIssuerSpecArgs{
			Acme: &certv1.ClusterIssuerSpecAcmeArgs{
				Email: pulumi.String(adminEmail),
				PrivateKeySecretRef: certv1.ClusterIssuerSpecAcmePrivateKeySecretRefArgs{
					Name: pulumi.String(clusterIssuerType.String()),
				},
				Server: pulumi.String(acmeServerURL),
				Solvers: &certv1.ClusterIssuerSpecAcmeSolversArray{
					&certv1.ClusterIssuerSpecAcmeSolversArgs{
						Http01: certv1.ClusterIssuerSpecAcmeSolversHttp01Args{
							Ingress: &certv1.ClusterIssuerSpecAcmeSolversHttp01IngressArgs{
								Class: pulumi.String(types.IngressClassNameTraefik.String()),
								IngressTemplate: certv1.ClusterIssuerSpecAcmeSolversHttp01IngressIngressTemplateArgs{
									Metadata: certv1.ClusterIssuerSpecAcmeSolversHttp01IngressIngressTemplateMetadataArgs{
										Annotations: pulumi.StringMap{
											//// needed for acmeCert with httpsRedirection
											"traefik.ingress.kubernetes.io/router.tls": pulumi.String("true"),
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}

	clusterIssuer, err := certv1.NewClusterIssuer(ctx, clusterIssuerType.String(), &certv1.ClusterIssuerArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(clusterIssuerType.String()),
		},
		Spec: clusterIssuerSpec,
	})
	if err != nil {
		return nil, err
	}

	return clusterIssuer, nil
}

func createCALocalSecret(ctx *pulumi.Context, clusterIssuerType types.ClusterIssuerType) error {
	conf := config.New(ctx, "")
	var ca types.CaSecret
	conf.RequireSecretObject("certManager", &ca)

	_, err := corev1.NewSecret(ctx, clusterIssuerType.String(), &corev1.SecretArgs{
		Data: pulumi.StringMap{
			"tls.crt": pulumi.String(ca.CA.Crt),
			"tls.key": pulumi.String(ca.CA.Key),
		},
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(clusterIssuerType.String()),
			Namespace: pulumi.String("cert-manager"),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func acmeServerURL(clusterIssuerType types.ClusterIssuerType) (string, error) {
	switch clusterIssuerType {
	case types.ClusterIssuerTypeCALocal:
		return "", errors.New("cannot assign acmeServerUrl to the ca-local-clusterissuer")
	case types.ClusterIssuerTypeLetsEncryptStaging:
		return "https://acme-staging-v02.api.letsencrypt.org/directory", nil
	case types.ClusterIssuerTypeLetsEncryptProd:
		return "https://acme-v02.api.letsencrypt.org/directory", nil
	}
	return "", errors.New("no acme-url specified")
}
