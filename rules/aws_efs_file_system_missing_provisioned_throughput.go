package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEFSFileSystemMissingProvisionedThroughputRule checks EFS provisioned throughput configuration.
type AwsEFSFileSystemMissingProvisionedThroughputRule struct {
	tflint.DefaultRule

	resourceType                   string
	throughputModeAttribute        string
	provisionedThroughputAttribute string
}

// NewAwsEFSFileSystemMissingProvisionedThroughputRule returns a new rule.
func NewAwsEFSFileSystemMissingProvisionedThroughputRule() *AwsEFSFileSystemMissingProvisionedThroughputRule {
	return &AwsEFSFileSystemMissingProvisionedThroughputRule{
		resourceType:                   "aws_efs_file_system",
		throughputModeAttribute:        "throughput_mode",
		provisionedThroughputAttribute: "provisioned_throughput_in_mibps",
	}
}

// Name returns the rule name.
func (r *AwsEFSFileSystemMissingProvisionedThroughputRule) Name() string {
	return "awscx_efs_file_system_missing_provisioned_throughput"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsEFSFileSystemMissingProvisionedThroughputRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsEFSFileSystemMissingProvisionedThroughputRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsEFSFileSystemMissingProvisionedThroughputRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system"
}

// Check reports when throughput_mode is provisioned without provisioned throughput.
func (r *AwsEFSFileSystemMissingProvisionedThroughputRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.throughputModeAttribute},
			{Name: r.provisionedThroughputAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		throughputMode, exists := resource.Body.Attributes[r.throughputModeAttribute]
		if !exists {
			continue
		}

		_, hasProvisionedThroughput := resource.Body.Attributes[r.provisionedThroughputAttribute]

		err := runner.EvaluateExpr(throughputMode.Expr, func(value string) error {
			if !strings.EqualFold(value, "provisioned") {
				return nil
			}
			if hasProvisionedThroughput {
				return nil
			}

			runner.EmitIssue(
				r,
				"`provisioned_throughput_in_mibps` must be set when `throughput_mode = \"provisioned\"`.",
				throughputMode.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
