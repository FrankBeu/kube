// *** WARNING: this file was generated by crd2pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v1

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Order is a type to represent an Order with an ACME server
type Order struct {
	pulumi.CustomResourceState

	ApiVersion pulumi.StringPtrOutput  `pulumi:"apiVersion"`
	Kind       pulumi.StringPtrOutput  `pulumi:"kind"`
	Metadata   metav1.ObjectMetaOutput `pulumi:"metadata"`
	Spec       OrderSpecOutput         `pulumi:"spec"`
	Status     OrderStatusPtrOutput    `pulumi:"status"`
}

// NewOrder registers a new resource with the given unique name, arguments, and options.
func NewOrder(ctx *pulumi.Context,
	name string, args *OrderArgs, opts ...pulumi.ResourceOption) (*Order, error) {
	if args == nil {
		args = &OrderArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("acme.cert-manager.io/v1")
	args.Kind = pulumi.StringPtr("Order")
	var resource Order
	err := ctx.RegisterResource("kubernetes:acme.cert-manager.io/v1:Order", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetOrder gets an existing Order resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetOrder(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *OrderState, opts ...pulumi.ResourceOption) (*Order, error) {
	var resource Order
	err := ctx.ReadResource("kubernetes:acme.cert-manager.io/v1:Order", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Order resources.
type orderState struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	Spec       *OrderSpec         `pulumi:"spec"`
	Status     *OrderStatus       `pulumi:"status"`
}

type OrderState struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	Spec       OrderSpecPtrInput
	Status     OrderStatusPtrInput
}

func (OrderState) ElementType() reflect.Type {
	return reflect.TypeOf((*orderState)(nil)).Elem()
}

type orderArgs struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	Spec       *OrderSpec         `pulumi:"spec"`
	Status     *OrderStatus       `pulumi:"status"`
}

// The set of arguments for constructing a Order resource.
type OrderArgs struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	Spec       OrderSpecPtrInput
	Status     OrderStatusPtrInput
}

func (OrderArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*orderArgs)(nil)).Elem()
}

type OrderInput interface {
	pulumi.Input

	ToOrderOutput() OrderOutput
	ToOrderOutputWithContext(ctx context.Context) OrderOutput
}

func (*Order) ElementType() reflect.Type {
	return reflect.TypeOf((*Order)(nil))
}

func (i *Order) ToOrderOutput() OrderOutput {
	return i.ToOrderOutputWithContext(context.Background())
}

func (i *Order) ToOrderOutputWithContext(ctx context.Context) OrderOutput {
	return pulumi.ToOutputWithContext(ctx, i).(OrderOutput)
}

type OrderOutput struct {
	*pulumi.OutputState
}

func (OrderOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Order)(nil))
}

func (o OrderOutput) ToOrderOutput() OrderOutput {
	return o
}

func (o OrderOutput) ToOrderOutputWithContext(ctx context.Context) OrderOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(OrderOutput{})
}
