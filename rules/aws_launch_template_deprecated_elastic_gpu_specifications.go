package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule warns on deprecated Elastic Graphics configuration.
type AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule returns a new rule.
func NewAwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule() *AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule {
	return &AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule{
		resourceType: "aws_launch_template",
		blockType:    "elastic_gpu_specifications",
	}
}

// Name returns the rule name.
func (r *AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule) Name() string {
	return "awscx_launch_template_deprecated_elastic_gpu_specifications"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template"
}

// Check warns when deprecated Elastic Graphics configuration is used on aws_launch_template.
func (r *AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule) Check(runner tflint.Runner) error {
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
				"`elastic_gpu_specifications` on `aws_launch_template` is deprecated because Amazon Elastic Graphics reached end of life on January 8, 2024.",
				block.DefRange,
			)
		}
	}

	return nil
}
