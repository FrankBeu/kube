// Package persistence  is responsible for creation of persistence-related resources lice pv and pvc
// on k3{d,s} a default local-path storageClass is already available
package persistence

import (
	"strings"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/kubeconfig"
)

//go:generate stringer -type=StorageClass -linecomment
const (
	StorageClassLocalPath StorageClass = iota // local-path
)

//go:generate stringer -type=AccessMode   -linecomment
const (
	AccessModeRWO AccessMode = iota // ReadWriteOnce
	AccessModeROX                   // ReadOnlyMany
	AccessModeRWX                   // ReadWriteMany
)
const (
	persistentVolumeAPIVersion      = "v1"
	persistentVolumeKind            = "PersistentVolume"
	persistentVolumeClaimAPIVersion = "v1"
	persistentVolumeClaimKind       = "PersistentVolumeClaim"

	persistentPathPrefix = "/srv/kube/persistence"
)

type StorageClass int

// TypeLabel returns the string value used as the typeLabel for on pv and pvc using according to the specified storageClass
func (s StorageClass) TypeLabel() string {
	strs := strings.Split(s.String(), "-")
	return strs[0]
}

type AccessMode int

type Config struct {
	Name         string
	StorageClass StorageClass
	Capacity     string
	AccessMode   AccessMode
	// relative to globalPrefix + env
	PersistencePathSuffix string
	NamespaceName         string
}

func CreatePersistentVolume(ctx *pulumi.Context, pc *Config) (*corev1.PersistentVolume, error) {
	pv, err := corev1.NewPersistentVolume(ctx, pc.Name, &corev1.PersistentVolumeArgs{
		ApiVersion: pulumi.String(persistentVolumeAPIVersion),
		Kind:       pulumi.String(persistentVolumeKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String(pc.Name),
			Labels: pulumi.StringMap{
				"type": pulumi.String(pc.StorageClass.TypeLabel()),
			},
		},
		Spec: &corev1.PersistentVolumeSpecArgs{
			StorageClassName: pulumi.String(pc.StorageClass.String()),
			Capacity: pulumi.StringMap{
				"storage": pulumi.String(pc.Capacity),
			},
			AccessModes: pulumi.StringArray{
				pulumi.String(pc.AccessMode.String()),
			},
			HostPath: &corev1.HostPathVolumeSourceArgs{
				Path: pulumi.String(persistentPathPrefix + "/" + kubeconfig.Env(ctx) + "/" + pc.PersistencePathSuffix),
			},
		},
	})
	return pv, err
}

func CreatePersistentVolumeClaim(ctx *pulumi.Context, pc *Config) (*corev1.PersistentVolumeClaim, error) {
	pvc, err := corev1.NewPersistentVolumeClaim(ctx, pc.Name, &corev1.PersistentVolumeClaimArgs{
		ApiVersion: pulumi.String(persistentVolumeClaimAPIVersion),
		Kind:       pulumi.String(persistentVolumeClaimKind),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(pc.Name),
			Namespace: pulumi.String(pc.NamespaceName),
		},
		Spec: &corev1.PersistentVolumeClaimSpecArgs{
			VolumeName:       pulumi.String(pc.Name),
			StorageClassName: pulumi.String(pc.StorageClass.String()),
			AccessModes: pulumi.StringArray{
				pulumi.String(pc.AccessMode.String()),
			},
			Resources: &corev1.ResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": pulumi.String(pc.Capacity),
				},
			},
		},
	})
	return pvc, err
}
