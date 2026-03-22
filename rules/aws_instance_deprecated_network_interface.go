package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsInstanceDeprecatedNetworkInterfaceRule warns when deprecated network_interface is used.
type AwsInstanceDeprecatedNetworkInterfaceRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsInstanceDeprecatedNetworkInterfaceRule returns a new rule.
func NewAwsInstanceDeprecatedNetworkInterfaceRule() *AwsInstanceDeprecatedNetworkInterfaceRule {
	return &AwsInstanceDeprecatedNetworkInterfaceRule{
		resourceType: "aws_instance",
		blockType:    "network_interface",
	}
}

// Name returns the rule name.
func (r *AwsInstanceDeprecatedNetworkInterfaceRule) Name() string {
	return "awscx_instance_deprecated_network_interface"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsInstanceDeprecatedNetworkInterfaceRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsInstanceDeprecatedNetworkInterfaceRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsInstanceDeprecatedNetworkInterfaceRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance"
}

// Check warns when deprecated network_interface blocks are configured on aws_instance.
func (r *AwsInstanceDeprecatedNetworkInterfaceRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: r.blockType},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			if block.Type != r.blockType {
				continue
			}

			runner.EmitIssue(
				r,
				"`network_interface` on `aws_instance` is deprecated; use `primary_network_interface` for the primary ENI and `aws_network_interface_attachment` for additional ENIs instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
