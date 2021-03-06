// *** WARNING: this file was generated by crd2pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v1

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Mapping is the Schema for the mappings API
type Mapping struct {
	pulumi.CustomResourceState

	ApiVersion pulumi.StringPtrOutput     `pulumi:"apiVersion"`
	Kind       pulumi.StringPtrOutput     `pulumi:"kind"`
	Metadata   metav1.ObjectMetaPtrOutput `pulumi:"metadata"`
	// MappingSpec defines the desired state of Mapping
	Spec MappingSpecPtrOutput `pulumi:"spec"`
	// MappingStatus defines the observed state of Mapping
	Status MappingStatusPtrOutput `pulumi:"status"`
}

// NewMapping registers a new resource with the given unique name, arguments, and options.
func NewMapping(ctx *pulumi.Context,
	name string, args *MappingArgs, opts ...pulumi.ResourceOption) (*Mapping, error) {
	if args == nil {
		args = &MappingArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("getambassador.io/v1")
	args.Kind = pulumi.StringPtr("Mapping")
	var resource Mapping
	err := ctx.RegisterResource("kubernetes:getambassador.io/v1:Mapping", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetMapping gets an existing Mapping resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetMapping(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *MappingState, opts ...pulumi.ResourceOption) (*Mapping, error) {
	var resource Mapping
	err := ctx.ReadResource("kubernetes:getambassador.io/v1:Mapping", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Mapping resources.
type mappingState struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// MappingSpec defines the desired state of Mapping
	Spec *MappingSpec `pulumi:"spec"`
	// MappingStatus defines the observed state of Mapping
	Status *MappingStatus `pulumi:"status"`
}

type MappingState struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// MappingSpec defines the desired state of Mapping
	Spec MappingSpecPtrInput
	// MappingStatus defines the observed state of Mapping
	Status MappingStatusPtrInput
}

func (MappingState) ElementType() reflect.Type {
	return reflect.TypeOf((*mappingState)(nil)).Elem()
}

type mappingArgs struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// MappingSpec defines the desired state of Mapping
	Spec *MappingSpec `pulumi:"spec"`
	// MappingStatus defines the observed state of Mapping
	Status *MappingStatus `pulumi:"status"`
}

// The set of arguments for constructing a Mapping resource.
type MappingArgs struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// MappingSpec defines the desired state of Mapping
	Spec MappingSpecPtrInput
	// MappingStatus defines the observed state of Mapping
	Status MappingStatusPtrInput
}

func (MappingArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*mappingArgs)(nil)).Elem()
}

type MappingInput interface {
	pulumi.Input

	ToMappingOutput() MappingOutput
	ToMappingOutputWithContext(ctx context.Context) MappingOutput
}

func (*Mapping) ElementType() reflect.Type {
	return reflect.TypeOf((*Mapping)(nil))
}

func (i *Mapping) ToMappingOutput() MappingOutput {
	return i.ToMappingOutputWithContext(context.Background())
}

func (i *Mapping) ToMappingOutputWithContext(ctx context.Context) MappingOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MappingOutput)
}

type MappingOutput struct {
	*pulumi.OutputState
}

func (MappingOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Mapping)(nil))
}

func (o MappingOutput) ToMappingOutput() MappingOutput {
	return o
}

func (o MappingOutput) ToMappingOutputWithContext(ctx context.Context) MappingOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(MappingOutput{})
}
