package types

type ClusterIssuerType int

//go:generate stringer -type=ClusterIssuerType -linecomment
const (
	ClusterIssuerTypeCALocal            ClusterIssuerType = iota // ca-local
	ClusterIssuerTypeLetsEncryptStaging                          // letsencrypt-staging
	ClusterIssuerTypeLetsEncryptProd                             // letsencrypt-prod

	DefaultDurationInDays = 90
)

var AllClusterIssuerTypes = []ClusterIssuerType{
	ClusterIssuerTypeCALocal,
	ClusterIssuerTypeLetsEncryptStaging,
	ClusterIssuerTypeLetsEncryptProd,
}

type Cert struct {
	ClusterIssuerType ClusterIssuerType
	Namespace         string
	// name is also used as subdomainName
	Name string
	// optional: default Duration is set to 90d
	Duration string
	// optional:
	AdditionalSubdomainNames []string
}

type CaSecret struct {
	CA struct {
		Crt string
		Key string
	}
}
