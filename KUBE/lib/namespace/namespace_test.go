package namespace

import (
	"sync"
	"testing"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"thesym.site/kube/lib/testutils"
)

func TestCreateNamespace(t *testing.T) {
	type args struct {
		nsInput *Namespace
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test a namespace with defaults",
			args: args{
				nsInput: &Namespace{
					Name: "testnamespace",
					Tier: NamespaceTierTesting,
				},
			},
		},
		{
			name: "test a namespace with glooDiscovery",
			args: args{
				nsInput: &Namespace{
					Name:          "gloo-discovery-testnamespace",
					Tier:          NamespaceTierEdge,
					GlooDiscovery: true,
				},
			},
		},
		{
			name: "test a namespace with additional labels",
			args: args{
				nsInput: &Namespace{
					Name: "additional-labels-testnamespace",
					Tier: NamespaceTierEdge,
					AdditionalLabels: []NamespaceLabel{
						{Name: "app.kubernetes.io/name", Value: "ingress-nginx"},
						{Name: "app.kubernetes.io/instance", Value: "ingress-nginx"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		//nolint:lll
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				nsResult, err := CreateNamespace(ctx, tt.args.nsInput)
				assert.NoError(t, err)

				var wg sync.WaitGroup
				wg.Add(1)

				pulumi.All(nsResult.ApiVersion, nsResult.Kind, nsResult.Metadata).ApplyT(func(all []interface{}) error {
					apivActual := *all[0].(*string)
					kindActual := *all[1].(*string)
					metaActual := *all[2].(*metav1.ObjectMeta)

					assert.Equal(t, apivActual, namespaceAPIVersion, "namespace has wrong apiVersion: wanted: %s, got: %s", namespaceAPIVersion, apivActual)
					assert.Equalf(t, kindActual, namespaceKind, "kind: %v is not a namespace", kindActual)
					assert.Equalf(t, *metaActual.Name, tt.args.nsInput.Name, "namespace has wrong metadata.Name: wanted: %s, got: %s", tt.args.nsInput.Name, *metaActual.Name)

					assertLabels(t, metaActual.Labels, tt.args.nsInput, tt.args.nsInput.GlooDiscovery)

					wg.Done()
					return nil
				})

				return nil
			}, pulumi.WithMocks("project", "stack", testutils.Mocks(0)))
			assert.NoError(t, err)
		})
	}
}

func assertLabels(t *testing.T, labelsActual testutils.Labels, ns *Namespace, glooDiscoveryActivated bool) {
	containsMsg := "label %q is not set"

	testutils.AssertLabel(t, labelsActual, "name", ns.Name, containsMsg)
	testutils.AssertLabel(t, labelsActual, "tier", ns.Tier.String(), containsMsg)

	assertGlooDiscovery(t, glooDiscoveryActivated, labelsActual)

	for _, label := range ns.AdditionalLabels {
		testutils.AssertLabel(t, labelsActual, label.Name, label.Value, containsMsg)
	}
}

func assertGlooDiscovery(t *testing.T, activated bool, labelsActual testutils.Labels) {
	containsMsgGloo := "glooDiscoveryLabel is NOT set despite glooDiscovery being activated"

	if activated {
		testutils.AssertLabel(t, labelsActual, namespaceGlooLabelKey, "enabled", containsMsgGloo)
	} else {
		assert.NotContainsf(
			t,
			labelsActual,
			namespaceGlooLabelKey,
			"glooDiscoveryLabel is set despite glooDiscovery NOT being activated",
		)
	}
}
