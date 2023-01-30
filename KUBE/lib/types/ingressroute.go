package types

const (
	IngressRouteAPIVersion = "traefik.containo.us/v1alpha1"
	IngressRouteKind       = "IngressRoute"
	// IngressRouteTLSSecretSuffix = "-ingressroute-tls" //nolint:gosec
	IngressRouteTLSSecretSuffix = "-tls"

	IngressRouteTLSDefaultDurationInDays = 99
)

type IngressRouteConfig struct {
	Name              string       // metadata.name; used also as simple 'match' and the 'cert.{,common}Name'; must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character
	NamespaceName     string       // metadata.namespace
	HostSubdomainName string       // routes[n].match - without the global domain
	Match             string       // routes[n].match if set replaces the HostSubdomainName in the match and the cert.{,common}Name
	Service           Service      // routes[n].services
	Middlewares       []Middleware //  routes[n].middlewares
}
type ServiceKind int

//go:generate stringer -type=ServiceKind -linecomment
const (
	ServiceKindK8S     ServiceKind = iota // Service
	ServiceKindTraefik                    // TraefikService
)

type Service struct {
	Kind          ServiceKind
	Name          string // routes[n].services[n].name
	NamespaceName string // routes[n].services[n].namespace
	Port          int    // routes[n].services[n].port
}

type MiddlewareName int

//go:generate stringer -type=MiddlewareName -linecomment
const (
	MiddleWareNameTest      MiddlewareName = iota // test
	MiddleWareNameBasicAuth                       // basicAuth

)

type Middleware struct {
	Name MiddlewareName // routes[n].middlewares[n].name
}
