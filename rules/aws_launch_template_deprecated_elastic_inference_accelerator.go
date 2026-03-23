package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule warns on deprecated Elastic Inference configuration.
type AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule returns a new rule.
func NewAwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule() *AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule {
	return &AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule{
		resourceType: "aws_launch_template",
		blockType:    "elastic_inference_accelerator",
	}
}

// Name returns the rule name.
func (r *AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule) Name() string {
	return "awscx_launch_template_deprecated_elastic_inference_accelerator"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule) Link() string {
	return "https://github.com/hashicorp/terraform-provider-aws/issues/41101"
}

// Check warns when deprecated Elastic Inference configuration is used on aws_launch_template.
func (r *AwsLaunchTemplateDeprecatedElasticInferenceAcceleratorRule) Check(runner tflint.Runner) error {
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
				"`elastic_inference_accelerator` on `aws_launch_template` is deprecated because Amazon Elastic Inference reached end of life in April 2024 and is no longer available.",
				block.DefRange,
			)
		}
	}

	return nil
}
