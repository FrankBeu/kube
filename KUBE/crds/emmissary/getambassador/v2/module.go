// *** WARNING: this file was generated by crd2pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package v2

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// A Module defines system-wide configuration.  The type of module is controlled by the .metadata.name; valid names are "ambassador" or "tls".
//  https://www.getambassador.io/docs/edge-stack/latest/topics/running/ambassador/#the-ambassador-module https://www.getambassador.io/docs/edge-stack/latest/topics/running/tls/#tls-module-deprecated
type Module struct {
	pulumi.CustomResourceState

	ApiVersion pulumi.StringPtrOutput     `pulumi:"apiVersion"`
	Kind       pulumi.StringPtrOutput     `pulumi:"kind"`
	Metadata   metav1.ObjectMetaPtrOutput `pulumi:"metadata"`
	Spec       ModuleSpecPtrOutput        `pulumi:"spec"`
}

// NewModule registers a new resource with the given unique name, arguments, and options.
func NewModule(ctx *pulumi.Context,
	name string, args *ModuleArgs, opts ...pulumi.ResourceOption) (*Module, error) {
	if args == nil {
		args = &ModuleArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("getambassador.io/v2")
	args.Kind = pulumi.StringPtr("Module")
	var resource Module
	err := ctx.RegisterResource("kubernetes:getambassador.io/v2:Module", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetModule gets an existing Module resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetModule(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ModuleState, opts ...pulumi.ResourceOption) (*Module, error) {
	var resource Module
	err := ctx.ReadResource("kubernetes:getambassador.io/v2:Module", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Module resources.
type moduleState struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	Spec       *ModuleSpec        `pulumi:"spec"`
}

type ModuleState struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	Spec       ModuleSpecPtrInput
}

func (ModuleState) ElementType() reflect.Type {
	return reflect.TypeOf((*moduleState)(nil)).Elem()
}

type moduleArgs struct {
	ApiVersion *string            `pulumi:"apiVersion"`
	Kind       *string            `pulumi:"kind"`
	Metadata   *metav1.ObjectMeta `pulumi:"metadata"`
	Spec       *ModuleSpec        `pulumi:"spec"`
}

// The set of arguments for constructing a Module resource.
type ModuleArgs struct {
	ApiVersion pulumi.StringPtrInput
	Kind       pulumi.StringPtrInput
	Metadata   metav1.ObjectMetaPtrInput
	Spec       ModuleSpecPtrInput
}

func (ModuleArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*moduleArgs)(nil)).Elem()
}

type ModuleInput interface {
	pulumi.Input

	ToModuleOutput() ModuleOutput
	ToModuleOutputWithContext(ctx context.Context) ModuleOutput
}

func (*Module) ElementType() reflect.Type {
	return reflect.TypeOf((*Module)(nil))
}

func (i *Module) ToModuleOutput() ModuleOutput {
	return i.ToModuleOutputWithContext(context.Background())
}

func (i *Module) ToModuleOutputWithContext(ctx context.Context) ModuleOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ModuleOutput)
}

type ModuleOutput struct {
	*pulumi.OutputState
}

func (ModuleOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Module)(nil))
}

func (o ModuleOutput) ToModuleOutput() ModuleOutput {
	return o
}

func (o ModuleOutput) ToModuleOutputWithContext(ctx context.Context) ModuleOutput {
	return o
}

func init() {
	pulumi.RegisterOutputType(ModuleOutput{})
}
