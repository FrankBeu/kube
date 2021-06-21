// *** WARNING: this file was generated by crd2pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v1

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A ClusterIssuer represents a certificate issuing authority which can be referenced as part of `issuerRef` fields. It is similar to an Issuer, however it is cluster-scoped and therefore can be referenced by resources that exist in *any* namespace, not just the same namespace as the referent.
type ClusterIssuer struct {
	pulumi.CustomResourceState

	ApiVersion pulumi.StringPtrOutput     `pulumi:"apiVersion"`
	Kind       pulumi.StringPtrOutput     `pulumi:"kind"`
	Metadata   metav1.ObjectMetaPtrOutput `pulumi:"metadata"`
	// Desired state of the ClusterIssuer resource.
	Spec ClusterIssuerSpecOutput `pulumi:"spec"`
	// Status of the ClusterIssuer. This is set and managed automatically.
	Status ClusterIssuerStatusPtrOutput `pulumi:"status"`
}

// NewClusterIssuer registers a new resource with the given unique name, arguments, and options.
func NewClusterIssuer(ctx *pulumi.Context,
	name string, args *ClusterIssuerArgs, opts ...pulumi.ResourceOption) (*ClusterIssuer, error) {
	if args == nil {
		args = &ClusterIssuerArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("cert-manager.io/v1")
	args.Kind = pulumi.StringPtr("ClusterIssuer")
	var resource ClusterIssuer
	err := ctx.RegisterResource("kubernetes:cert-manager.io/v1:ClusterIssuer", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetClusterIssuer gets an existing ClusterIssuer resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetClusterIssuer(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ClusterIssuerState, opts ...pulumi.ResourceOption) (*ClusterIssuer, error) {
	var resource ClusterIssuer
	err := ctx.ReadResource("kubernetes:cert-manager.io/v1:ClusterIssuer", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ClusterIssuer resources.
type clusterIssuerState struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// Desired state of the ClusterIssuer resource.
	Spec *ClusterIssuerSpec `pulumi:"spec"`
	// Status of the ClusterIssuer. This is set and managed automatically.
	Status *ClusterIssuerStatus `pulumi:"status"`
}

type ClusterIssuerState struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// Desired state of the ClusterIssuer resource.
	Spec ClusterIssuerSpecPtrInput
	// Status of the ClusterIssuer. This is set and managed automatically.
	Status ClusterIssuerStatusPtrInput
}

func (ClusterIssuerState) ElementType() reflect.Type {
	return reflect.TypeOf((*clusterIssuerState)(nil)).Elem()
}

type clusterIssuerArgs struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// Desired state of the ClusterIssuer resource.
	Spec *ClusterIssuerSpec `pulumi:"spec"`
	// Status of the ClusterIssuer. This is set and managed automatically.
	Status *ClusterIssuerStatus `pulumi:"status"`
}

// The set of arguments for constructing a ClusterIssuer resource.
type ClusterIssuerArgs struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// Desired state of the ClusterIssuer resource.
	Spec ClusterIssuerSpecPtrInput
	// Status of the ClusterIssuer. This is set and managed automatically.
	Status ClusterIssuerStatusPtrInput
}

func (ClusterIssuerArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*clusterIssuerArgs)(nil)).Elem()
}

type ClusterIssuerInput interface {
	pulumi.Input

	ToClusterIssuerOutput() ClusterIssuerOutput
	ToClusterIssuerOutputWithContext(ctx context.Context) ClusterIssuerOutput
}

func (*ClusterIssuer) ElementType() reflect.Type {
	return reflect.TypeOf((*ClusterIssuer)(nil))
}

func (i *ClusterIssuer) ToClusterIssuerOutput() ClusterIssuerOutput {
	return i.ToClusterIssuerOutputWithContext(context.Background())
}

func (i *ClusterIssuer) ToClusterIssuerOutputWithContext(ctx context.Context) ClusterIssuerOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ClusterIssuerOutput)
}

type ClusterIssuerOutput struct {
	*pulumi.OutputState
}

func (ClusterIssuerOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*ClusterIssuer)(nil))
}

func (o ClusterIssuerOutput) ToClusterIssuerOutput() ClusterIssuerOutput {
	return o
}

func (o ClusterIssuerOutput) ToClusterIssuerOutputWithContext(ctx context.Context) ClusterIssuerOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(ClusterIssuerOutput{})
}
