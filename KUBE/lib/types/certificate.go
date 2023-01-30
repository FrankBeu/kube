package types

type ClusterIssuerType int

//go:generate stringer -type=ClusterIssuerType -linecomment
const (
	ClusterIssuerTypeCALocal            ClusterIssuerType = iota // ca-local
	ClusterIssuerTypeLetsEncryptStaging                          // letsencrypt-staging
	ClusterIssuerTypeLetsEncryptProd                             // letsencrypt-prod

	CertificateDefaultDurationInDays = 99
)

var AllClusterIssuerTypes = []ClusterIssuerType{
	ClusterIssuerTypeCALocal,
	ClusterIssuerTypeLetsEncryptStaging,
	ClusterIssuerTypeLetsEncryptProd,
}

type Cert struct {
	ClusterIssuerType        ClusterIssuerType
	Namespace                string
	Name                     string   // metadata.name
	Duration                 string   // optional: spec.duration; default Duration is set to 90d
	CommonNameSegment        string   // spec.commonName - subdomain (the domain is concated automatically)
	AdditionalSubdomainNames []string // optional
}

type CaSecret struct {
	CA struct {
		Crt string
		Key string
	}
}
