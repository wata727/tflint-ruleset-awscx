package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEFSFileSystemProvisionedThroughputNonProvisionedRule checks that provisioned throughput is only used with provisioned mode.
type AwsEFSFileSystemProvisionedThroughputNonProvisionedRule struct {
	tflint.DefaultRule

	resourceType                   string
	throughputModeAttribute        string
	provisionedThroughputAttribute string
}

// NewAwsEFSFileSystemProvisionedThroughputNonProvisionedRule returns a new rule.
func NewAwsEFSFileSystemProvisionedThroughputNonProvisionedRule() *AwsEFSFileSystemProvisionedThroughputNonProvisionedRule {
	return &AwsEFSFileSystemProvisionedThroughputNonProvisionedRule{
		resourceType:                   "aws_efs_file_system",
		throughputModeAttribute:        "throughput_mode",
		provisionedThroughputAttribute: "provisioned_throughput_in_mibps",
	}
}

// Name returns the rule name.
func (r *AwsEFSFileSystemProvisionedThroughputNonProvisionedRule) Name() string {
	return "awscx_efs_file_system_provisioned_throughput_non_provisioned"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsEFSFileSystemProvisionedThroughputNonProvisionedRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsEFSFileSystemProvisionedThroughputNonProvisionedRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsEFSFileSystemProvisionedThroughputNonProvisionedRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system"
}

// Check reports when provisioned throughput is configured without provisioned throughput mode.
func (r *AwsEFSFileSystemProvisionedThroughputNonProvisionedRule) Check(runner tflint.Runner) error {
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
		provisionedThroughput, hasProvisionedThroughput := resource.Body.Attributes[r.provisionedThroughputAttribute]
		if !hasProvisionedThroughput {
			continue
		}

		throughputMode, hasThroughputMode := resource.Body.Attributes[r.throughputModeAttribute]
		if !hasThroughputMode {
			runner.EmitIssue(
				r,
				"`provisioned_throughput_in_mibps` can only be set when `throughput_mode = \"provisioned\"`.",
				provisionedThroughput.Expr.Range(),
			)
			continue
		}

		err := runner.EvaluateExpr(throughputMode.Expr, func(value string) error {
			if strings.EqualFold(value, "provisioned") {
				return nil
			}

			runner.EmitIssue(
				r,
				"`provisioned_throughput_in_mibps` can only be set when `throughput_mode = \"provisioned\"`.",
				provisionedThroughput.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
