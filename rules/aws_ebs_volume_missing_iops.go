package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEBSVolumeMissingIOPSRule checks provisioned IOPS volume configuration.
type AwsEBSVolumeMissingIOPSRule struct {
	tflint.DefaultRule

	resourceType  string
	typeAttribute string
	iopsAttribute string
}

// NewAwsEBSVolumeMissingIOPSRule returns a new rule.
func NewAwsEBSVolumeMissingIOPSRule() *AwsEBSVolumeMissingIOPSRule {
	return &AwsEBSVolumeMissingIOPSRule{
		resourceType:  "aws_ebs_volume",
		typeAttribute: "type",
		iopsAttribute: "iops",
	}
}

// Name returns the rule name.
func (r *AwsEBSVolumeMissingIOPSRule) Name() string {
	return "awscx_ebs_volume_missing_iops"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsEBSVolumeMissingIOPSRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsEBSVolumeMissingIOPSRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsEBSVolumeMissingIOPSRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume"
}

// Check reports when io1/io2 volume types omit iops.
func (r *AwsEBSVolumeMissingIOPSRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.typeAttribute},
			{Name: r.iopsAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		typeAttr, exists := resource.Body.Attributes[r.typeAttribute]
		if !exists {
			continue
		}

		_, hasIOPS := resource.Body.Attributes[r.iopsAttribute]

		err := runner.EvaluateExpr(typeAttr.Expr, func(value string) error {
			normalized := strings.ToLower(value)
			if normalized != "io1" && normalized != "io2" {
				return nil
			}
			if hasIOPS {
				return nil
			}

			runner.EmitIssue(
				r,
				"`iops` must be set when `type` is `io1` or `io2`.",
				typeAttr.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
