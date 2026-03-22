package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEBSVolumeThroughputNonGP3Rule checks that throughput is only used with gp3 volumes.
type AwsEBSVolumeThroughputNonGP3Rule struct {
	tflint.DefaultRule

	resourceType        string
	typeAttribute       string
	throughputAttribute string
}

// NewAwsEBSVolumeThroughputNonGP3Rule returns a new rule.
func NewAwsEBSVolumeThroughputNonGP3Rule() *AwsEBSVolumeThroughputNonGP3Rule {
	return &AwsEBSVolumeThroughputNonGP3Rule{
		resourceType:        "aws_ebs_volume",
		typeAttribute:       "type",
		throughputAttribute: "throughput",
	}
}

// Name returns the rule name.
func (r *AwsEBSVolumeThroughputNonGP3Rule) Name() string {
	return "awscx_ebs_volume_throughput_non_gp3"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsEBSVolumeThroughputNonGP3Rule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsEBSVolumeThroughputNonGP3Rule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsEBSVolumeThroughputNonGP3Rule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume"
}

// Check reports when throughput is configured without gp3 volume type.
func (r *AwsEBSVolumeThroughputNonGP3Rule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.typeAttribute},
			{Name: r.throughputAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		throughput, hasThroughput := resource.Body.Attributes[r.throughputAttribute]
		if !hasThroughput {
			continue
		}

		volumeType, hasType := resource.Body.Attributes[r.typeAttribute]
		if !hasType {
			runner.EmitIssue(
				r,
				"`throughput` can only be set when `type = \"gp3\"`.",
				throughput.Expr.Range(),
			)
			continue
		}

		err := runner.EvaluateExpr(volumeType.Expr, func(value string) error {
			if strings.EqualFold(value, "gp3") {
				return nil
			}

			runner.EmitIssue(
				r,
				"`throughput` can only be set when `type = \"gp3\"`.",
				throughput.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
