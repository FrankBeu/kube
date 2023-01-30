package persistence

import (
	"sync"
	"testing"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/stretchr/testify/assert"

	"thesym.site/kube/lib/kubeconfig"
	"thesym.site/kube/lib/testutil"
)

type args struct {
	persistenceConfig Config
	testConfig        map[string]string
}

var tests = []struct {
	name string
	args args
}{
	{
		name: "default pv",
		args: args{
			persistenceConfig: Config{
				Name:                  "test",
				StorageClass:          StorageClassLocalPath,
				Capacity:              "20Gi",
				AccessMode:            AccessModeROX,
				PersistencePathSuffix: "test/db",
				NamespaceName:         "test",
			},
			testConfig: map[string]string{
				"project:env": "dev",
			},
		},
	},
	{
		name: "giteaData claim",
		args: args{
			persistenceConfig: Config{
				Name:                  "gitea-data",
				StorageClass:          StorageClassLocalPath,
				Capacity:              "10Gi",
				AccessMode:            AccessModeRWO,
				PersistencePathSuffix: "gitea/data",
				NamespaceName:         "gitea",
			},
			testConfig: map[string]string{
				"project:env": "prod",
			},
		},
	},
	{
		name: "giteaDb claim",
		args: args{
			persistenceConfig: Config{
				Name:                  "data-gitea-postgresql-0",
				StorageClass:          StorageClassLocalPath,
				Capacity:              "10Gi",
				AccessMode:            AccessModeRWO,
				PersistencePathSuffix: "gitea/db",
				NamespaceName:         "gitea",
			},
			testConfig: map[string]string{
				"project:env": "prod",
			},
		},
	},
}

func TestCreatePersistentVolume(t *testing.T) {
	for _, tt := range tests {
		//nolint:lll
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				pvActual, err := CreatePersistentVolume(ctx, &tt.args.persistenceConfig)
				assert.NoError(t, err)

				var wg sync.WaitGroup
				wg.Add(1)

				pulumi.All(pvActual.ApiVersion, pvActual.Kind, pvActual.Metadata, pvActual.Spec).ApplyT(func(all []interface{}) error {
					apivActual := *all[0].(*string)
					kindActual := *all[1].(*string)
					metaActual := *all[2].(*metav1.ObjectMeta)
					specActual := *all[3].(*corev1.PersistentVolumeSpec)

					testutil.Equalf(t, "PersistentVolume", "ApiVersion", apivActual, persistentVolumeAPIVersion) //// v1 is always set
					testutil.Equalf(t, "PersistentVolume", "Kind", kindActual, persistentVolumeKind)             //// PersistentVolume is always set
					testutil.Equalf(t, "PersistentVolume", "Name", *metaActual.Name, tt.args.persistenceConfig.Name)
					testutil.AssertLabel(t, metaActual.Labels, "type", tt.args.persistenceConfig.StorageClass.TypeLabel(), "label %q is not set")
					testutil.Equalf(t, "PersistentVolume", "Spec.StorageClass", *specActual.StorageClassName, tt.args.persistenceConfig.StorageClass.String())
					testutil.Equalf(t, "PersistentVolume", "Spec.StorageCapacity", specActual.Capacity["storage"], tt.args.persistenceConfig.Capacity)
					testutil.Containsf(t, "PersistentVolume", "Spec.AccessMode", specActual.AccessModes, tt.args.persistenceConfig.AccessMode.String())

					env := kubeconfig.Env(ctx)
					persistencePathTarget := persistentPathPrefix + "/" + env + "/" + tt.args.persistenceConfig.PersistencePathSuffix
					testutil.Equalf(t, "PersistentVolume", "Spec.HostPath.Path/PersistencePath", specActual.HostPath.Path, persistencePathTarget)

					wg.Done()
					return nil
				})

				return nil
			}, testutil.WithMocksAndConfig("project", "stack", tt.args.testConfig, testutil.Mocks(0)))
			assert.NoError(t, err)
		})
	}
}

func TestCreatePersistentVolumeClaim(t *testing.T) {
	for _, tt := range tests {
		//nolint:lll
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				pvcActual, err := CreatePersistentVolumeClaim(ctx, &tt.args.persistenceConfig)
				assert.NoError(t, err)

				var wg sync.WaitGroup
				wg.Add(1)

				pulumi.All(pvcActual.ApiVersion, pvcActual.Kind, pvcActual.Metadata, pvcActual.Spec).ApplyT(func(all []interface{}) error {
					apivActual := *all[0].(*string)
					kindActual := *all[1].(*string)
					metaActual := *all[2].(*metav1.ObjectMeta)
					specActual := *all[3].(*corev1.PersistentVolumeClaimSpec)

					testutil.Equalf(t, "PersistentVolumeClaim", "ApiVersion", apivActual, persistentVolumeClaimAPIVersion) //// v1 is always set
					testutil.Equalf(t, "PersistentVolumeClaim", "Kind", kindActual, persistentVolumeClaimKind)             //// PersistentVolume is always set
					testutil.Equalf(t, "PersistentVolumeClaim", "Name", *metaActual.Name, tt.args.persistenceConfig.Name)
					testutil.Equalf(t, "PersistentVolumeClaim", "Namespace", *metaActual.Namespace, tt.args.persistenceConfig.NamespaceName)
					testutil.Equalf(t, "PersistentVolumeClaim", "Spec.VolumeName", *specActual.VolumeName, tt.args.persistenceConfig.Name)
					testutil.Equalf(t, "PersistentVolumeClaim", "Spec.StorageClass", *specActual.StorageClassName, tt.args.persistenceConfig.StorageClass.String())
					testutil.Equalf(t, "PersistentVolumeClaim", "Spec.Resources.Requiests.storage", specActual.Resources.Requests["storage"], tt.args.persistenceConfig.Capacity)
					testutil.Containsf(t, "PersistentVolumeClaim", "Spec.AccessMode", specActual.AccessModes, tt.args.persistenceConfig.AccessMode.String())

					wg.Done()
					return nil
				})

				return nil
			}, pulumi.WithMocks("project", "stack", testutil.Mocks(0)))
			assert.NoError(t, err)
		})
	}
}

func TestStorageClass_TypeLabel(t *testing.T) {
	tests := []struct {
		name string
		s    StorageClass
		want string
	}{
		{
			name: "test the currently used storageClass",
			s:    StorageClassLocalPath,
			want: "local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.TypeLabel(); got != tt.want {
				t.Errorf("StorageClass.TypeLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}
