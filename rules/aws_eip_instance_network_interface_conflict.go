package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEIPInstanceNetworkInterfaceConflictRule checks invalid aws_eip association arguments.
type AwsEIPInstanceNetworkInterfaceConflictRule struct {
	tflint.DefaultRule

	resourceType              string
	instanceAttribute         string
	networkInterfaceAttribute string
}

// NewAwsEIPInstanceNetworkInterfaceConflictRule returns a new rule.
func NewAwsEIPInstanceNetworkInterfaceConflictRule() *AwsEIPInstanceNetworkInterfaceConflictRule {
	return &AwsEIPInstanceNetworkInterfaceConflictRule{
		resourceType:              "aws_eip",
		instanceAttribute:         "instance",
		networkInterfaceAttribute: "network_interface",
	}
}

// Name returns the rule name.
func (r *AwsEIPInstanceNetworkInterfaceConflictRule) Name() string {
	return "awscx_eip_instance_network_interface_conflict"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsEIPInstanceNetworkInterfaceConflictRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsEIPInstanceNetworkInterfaceConflictRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsEIPInstanceNetworkInterfaceConflictRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip"
}

// Check reports when aws_eip sets both instance and network_interface.
func (r *AwsEIPInstanceNetworkInterfaceConflictRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.instanceAttribute},
			{Name: r.networkInterfaceAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		_, hasInstance := resource.Body.Attributes[r.instanceAttribute]
		networkInterface, hasNetworkInterface := resource.Body.Attributes[r.networkInterfaceAttribute]
		if !hasInstance || !hasNetworkInterface {
			continue
		}

		runner.EmitIssue(
			r,
			"`instance` and `network_interface` on `aws_eip` are mutually exclusive; set only one association target.",
			networkInterface.Expr.Range(),
		)
	}

	return nil
}
