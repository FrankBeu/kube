package lib

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	certv1 "thesym.site/kube/crds/cert-manager/certmanager/v1"
	"strconv"
)

type ClusterIssuerType int

//go:generate stringer -type=ClusterIssuerType -linecomment
const (
	ClusterIssuerTypeCaLocal            ClusterIssuerType = iota // ca-local
	ClusterIssuerTypeLetsEncryptProd                             // letsencrypt-prod
	ClusterIssuerTypeLetsEncryptStaging                          // letsencrypt-staging
)

type Cert struct {
	ClusterIssuerType  ClusterIssuerType
	Namespace          string
	// name is also used as subdomainName
	Name               string
	// optional: default Duration is set to 90d
	Duration           string
	// optional:
	AdditionalSubdomainNames []string
}

func CreateCert(ctx *pulumi.Context, cert *Cert) error {
	//// defaultDuration: 90 days
	if cert.Duration == "" {
		cert.Duration = strconv.Itoa(90*24) + "h"
	}

	conf := config.New(ctx, "")
	domainNameSuffix := "." + conf.Require("domain")
	domainName := cert.Name + domainNameSuffix


	dnsNames  := pulumi.StringArray{
		pulumi.String(domainName),
	}
	for _, subdomainName := range cert.AdditionalSubdomainNames {
		dnsNames = append(dnsNames, pulumi.String(subdomainName + domainNameSuffix))
	}

	_, err := certv1.NewCertificate(ctx, cert.Name, &certv1.CertificateArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(cert.Name),
			Namespace: pulumi.String(cert.Namespace),
		},
		Spec: certv1.CertificateSpecArgs{
			SecretName:  pulumi.String(cert.Name + "-tls"),
			Duration:    pulumi.String(cert.Duration),
			RenewBefore: pulumi.String("360h"),
			CommonName:  pulumi.String(domainName),
			DnsNames: dnsNames,
			IssuerRef: certv1.CertificateSpecIssuerRefArgs{
				Name: pulumi.String(cert.ClusterIssuerType.String()),
				//// currently only ClusterIssuers are used
				Kind: pulumi.String("ClusterIssuer"),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

type CaSecret struct {
	CA struct {
		Crt string
		Key string
	}
}

// createClusterIssuer creates an internal clusterIssuer with a local CA or an letsencrypt based issuer
func CreateClusterIssuer(ctx *pulumi.Context, clusterIssuerType ClusterIssuerType) error {
	clusterIssuerSpec := &certv1.ClusterIssuerSpecArgs{}

	if clusterIssuerType == ClusterIssuerTypeCaLocal {
		clusterIssuerSecretName := clusterIssuerType.String()
		conf := config.New(ctx, "")
		var ca CaSecret
		conf.RequireSecretObject("certManager", &ca)

		_, err := corev1.NewSecret(ctx, clusterIssuerType.String(), &corev1.SecretArgs{
			Data: pulumi.StringMap{
				"tls.crt": pulumi.String(ca.CA.Crt),
				"tls.key": pulumi.String(ca.CA.Key),
			},
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(clusterIssuerSecretName),
				Namespace: pulumi.String("cert-manager"),
			},
		})
		if err != nil {
			return err
		}

		clusterIssuerSpec = &certv1.ClusterIssuerSpecArgs{
			Ca: &certv1.ClusterIssuerSpecCaArgs{
				// CrlDistributionPoints: []string{},
				// OcspServers:           []string{},
				SecretName: pulumi.String(clusterIssuerSecretName),
			},
		}
	} else if clusterIssuerType == ClusterIssuerTypeLetsEncryptProd {
		clusterIssuerSpec = &certv1.ClusterIssuerSpecArgs{
			Acme: &certv1.ClusterIssuerSpecAcmeArgs{

				// TODO:
				// Status: nil,
				// apiVersion: cert-manager.io/v1
				// kind: ClusterIssuer
				// metadata:
				//   name: letsencrypt-prod
				// spec:
				//   acme:
				//     email: {{.hostAcmeEmail}}
				//     server: https://acme-v02.api.letsencrypt.org/directory
				//     privateKeySecretRef:
				//       name: letsencrypt-prod
				//     solvers:
				//     - http01:
				//         ingress:
				//           class: nginx
				//       selector: {}
			},
		}
	} else if clusterIssuerType == ClusterIssuerTypeLetsEncryptStaging {

		//////////////////////////////////////////////////////////////
		// hostApply:
		//   desc: apply the host for {{.name}}
		//   cmds:
		//     - |
		//       cat <<EOF | kubectl apply -f -
		//       apiVersion: getambassador.io/v2
		//       kind: Host
		//       metadata:
		//         name: {{.name}}
		//         namespace: ambassador-hosts
		//       spec:
		//         hostname: {{.domainName}}
		//         acmeProvider:
		//           {{- if .hostCertStaging }}
		//           authority: https://acme-staging-v02.api.letsencrypt.org/directory
		//           {{- else}}
		//           authority: https://acme-v02.api.letsencrypt.org/directory
		//           {{- end}}
		//           email: {{.hostAcmeEmail}}
		//           privateKeySecret:
		//             name: {{.name}}{{- if .hostCertStaging }}Staging{{- end}}-key
		//         tlsSecret:
		//           name: {{.name}}{{- if .hostCertStaging }}Staging{{- end}}
		//         requestPolicy:
		//           insecure:
		//             action: Redirect
	}

	_, err := certv1.NewClusterIssuer(ctx, clusterIssuerType.String(), &certv1.ClusterIssuerArgs{
		// ApiVersion: nil,
		// Kind:       nil,
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(clusterIssuerType.String()),
		},
		Spec: clusterIssuerSpec,
	})
	if err != nil {
		return err
	}

	return nil
}
