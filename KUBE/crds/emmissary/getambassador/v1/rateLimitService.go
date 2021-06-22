// *** WARNING: this file was generated by crd2pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v1

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// RateLimitService is the Schema for the ratelimitservices API
type RateLimitService struct {
	pulumi.CustomResourceState

	ApiVersion pulumi.StringPtrOutput     `pulumi:"apiVersion"`
	Kind       pulumi.StringPtrOutput     `pulumi:"kind"`
	Metadata   metav1.ObjectMetaPtrOutput `pulumi:"metadata"`
	// RateLimitServiceSpec defines the desired state of RateLimitService
	Spec RateLimitServiceSpecPtrOutput `pulumi:"spec"`
}

// NewRateLimitService registers a new resource with the given unique name, arguments, and options.
func NewRateLimitService(ctx *pulumi.Context,
	name string, args *RateLimitServiceArgs, opts ...pulumi.ResourceOption) (*RateLimitService, error) {
	if args == nil {
		args = &RateLimitServiceArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("getambassador.io/v1")
	args.Kind = pulumi.StringPtr("RateLimitService")
	var resource RateLimitService
	err := ctx.RegisterResource("kubernetes:getambassador.io/v1:RateLimitService", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetRateLimitService gets an existing RateLimitService resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetRateLimitService(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *RateLimitServiceState, opts ...pulumi.ResourceOption) (*RateLimitService, error) {
	var resource RateLimitService
	err := ctx.ReadResource("kubernetes:getambassador.io/v1:RateLimitService", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering RateLimitService resources.
type rateLimitServiceState struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// RateLimitServiceSpec defines the desired state of RateLimitService
	Spec *RateLimitServiceSpec `pulumi:"spec"`
}

type RateLimitServiceState struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// RateLimitServiceSpec defines the desired state of RateLimitService
	Spec RateLimitServiceSpecPtrInput
}

func (RateLimitServiceState) ElementType() reflect.Type {
	return reflect.TypeOf((*rateLimitServiceState)(nil)).Elem()
}

type rateLimitServiceArgs struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	// RateLimitServiceSpec defines the desired state of RateLimitService
	Spec *RateLimitServiceSpec `pulumi:"spec"`
}

// The set of arguments for constructing a RateLimitService resource.
type RateLimitServiceArgs struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	// RateLimitServiceSpec defines the desired state of RateLimitService
	Spec RateLimitServiceSpecPtrInput
}

func (RateLimitServiceArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*rateLimitServiceArgs)(nil)).Elem()
}

type RateLimitServiceInput interface {
	pulumi.Input

	ToRateLimitServiceOutput() RateLimitServiceOutput
	ToRateLimitServiceOutputWithContext(ctx context.Context) RateLimitServiceOutput
}

func (*RateLimitService) ElementType() reflect.Type {
	return reflect.TypeOf((*RateLimitService)(nil))
}

func (i *RateLimitService) ToRateLimitServiceOutput() RateLimitServiceOutput {
	return i.ToRateLimitServiceOutputWithContext(context.Background())
}

func (i *RateLimitService) ToRateLimitServiceOutputWithContext(ctx context.Context) RateLimitServiceOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RateLimitServiceOutput)
}

type RateLimitServiceOutput struct {
	*pulumi.OutputState
}

func (RateLimitServiceOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*RateLimitService)(nil))
}

func (o RateLimitServiceOutput) ToRateLimitServiceOutput() RateLimitServiceOutput {
	return o
}

func (o RateLimitServiceOutput) ToRateLimitServiceOutputWithContext(ctx context.Context) RateLimitServiceOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(RateLimitServiceOutput{})
}
