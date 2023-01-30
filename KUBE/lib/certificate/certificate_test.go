package certificate

import (
	"sync"
	"testing"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	certv1 "thesym.site/kube/crds/kubernetes/certmanager/v1"
	"thesym.site/kube/lib/testutil"
	"thesym.site/kube/lib/types"
)

// TODO write tests
func TestCreateCert(t *testing.T) {
	type args struct {
		ctx  *pulumi.Context
		cert *types.Cert
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CreateCert(tt.args.ctx, tt.args.cert); (err != nil) != tt.wantErr {
				t.Errorf("CreateCert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createClusterIssuer(t *testing.T) {
	//// TODO test exported function with secrets
	////      there is currently no cheap way of testing with applied secrets - only unwrapped function is tested
	////      XOR:
	////          Find way
	////          wait for https://github.com/pulumi/pulumi/issues/4472 to include secrets??
	type args struct {
		clusterIssuerType types.ClusterIssuerType
		adminEmail        string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: types.ClusterIssuerTypeCALocal.String(),
			args: args{
				clusterIssuerType: types.ClusterIssuerTypeCALocal,
				adminEmail:        "testCALocal@email",
			},
		},
		{

			name: types.ClusterIssuerTypeLetsEncryptStaging.String(),
			args: args{
				clusterIssuerType: types.ClusterIssuerTypeLetsEncryptStaging,
				adminEmail:        "testStaging@email",
			},
		},
		{
			name: types.ClusterIssuerTypeLetsEncryptProd.String(),
			args: args{
				clusterIssuerType: types.ClusterIssuerTypeLetsEncryptProd,
				adminEmail:        "testProduction@email",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				clusterIssuerResult, err := createClusterIssuer(ctx, tt.args.clusterIssuerType, tt.args.adminEmail)
				assert.NoError(t, err)

				var wg sync.WaitGroup
				wg.Add(1)

				//nolint:lll
				pulumi.All(clusterIssuerResult.ApiVersion, clusterIssuerResult.Kind, clusterIssuerResult.Metadata, clusterIssuerResult.Spec).ApplyT(func(all []interface{}) error {
					apivActual := *all[0].(*string)
					kindActual := *all[1].(*string)
					metaActual := *all[2].(*metav1.ObjectMeta)
					specActual := all[3].(certv1.ClusterIssuerSpec)

					testutil.Equalf(t, "ClusterIssuer", "ApiVersion", apivActual, "cert-manager.io/v1")
					testutil.Equalf(t, "ClusterIssuer", "Kind", kindActual, "ClusterIssuer")
					testutil.Equalf(t, "ClusterIssuer", "Name", *metaActual.Name, tt.args.clusterIssuerType.String())

					if tt.args.clusterIssuerType.String() == types.ClusterIssuerTypeCALocal.String() {
						testutil.Equalf(t, "ClusterIssuer-caLocal", "SecretName", specActual.Ca.SecretName, tt.args.clusterIssuerType.String())
					}

					if ct := tt.args.clusterIssuerType.String(); ct == types.ClusterIssuerTypeLetsEncryptStaging.String() || ct == types.ClusterIssuerTypeLetsEncryptProd.String() {
						testutil.Equalf(t, "ClusterIssuer-ACME", "email", *specActual.Acme.Email, tt.args.adminEmail)
						testutil.Equalf(t, "ClusterIssuer-ACME", "PrivateKeySecretRef", specActual.Acme.PrivateKeySecretRef.Name, tt.args.clusterIssuerType.String())

						acmeServerURL, err := acmeServerURL(tt.args.clusterIssuerType)
						assert.NoError(t, err)
						testutil.Equalf(t, "ClusterIssuer-ACME", "acmeServerUrl", specActual.Acme.Server, acmeServerURL)

						//// currently only one Solver is used per clusterIssuer
						testutil.Equalf(t, "ClusterIssuer-ACME", "solverIngressClass", *specActual.Acme.Solvers[0].Http01.Ingress.Class, types.IngressClassNameTraefik.String())
					}

					wg.Done()
					return nil
				})

				return nil
			}, testutil.WithMocksAndConfig("project", "stack", testutil.TestConfig, testutil.Mocks(0)))
			assert.NoError(t, err)
		})
	}
}
