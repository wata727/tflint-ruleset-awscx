package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLaunchTemplateIMDSv2OptionalTokensRule warns when launch templates explicitly allow IMDSv1.
type AwsLaunchTemplateIMDSv2OptionalTokensRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
	attribute    string
}

// NewAwsLaunchTemplateIMDSv2OptionalTokensRule returns a new rule.
func NewAwsLaunchTemplateIMDSv2OptionalTokensRule() *AwsLaunchTemplateIMDSv2OptionalTokensRule {
	return &AwsLaunchTemplateIMDSv2OptionalTokensRule{
		resourceType: "aws_launch_template",
		blockType:    "metadata_options",
		attribute:    "http_tokens",
	}
}

// Name returns the rule name.
func (r *AwsLaunchTemplateIMDSv2OptionalTokensRule) Name() string {
	return "awscx_launch_template_imdsv2_optional_tokens"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLaunchTemplateIMDSv2OptionalTokensRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLaunchTemplateIMDSv2OptionalTokensRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsLaunchTemplateIMDSv2OptionalTokensRule) Link() string {
	return "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html"
}

// Check warns when launch templates explicitly configure IMDSv1 compatibility.
func (r *AwsLaunchTemplateIMDSv2OptionalTokensRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: r.attribute}},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			attribute, exists := block.Body.Attributes[r.attribute]
			if !exists {
				continue
			}

			err := runner.EvaluateExpr(attribute.Expr, func(value string) error {
				if !strings.EqualFold(value, "optional") {
					return nil
				}

				runner.EmitIssue(
					r,
					"`metadata_options.http_tokens = \"optional\"` allows IMDSv1; prefer `required` to enforce IMDSv2.",
					attribute.Expr.Range(),
				)
				return nil
			}, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
