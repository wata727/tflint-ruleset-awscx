package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsAPIGatewayDeploymentDeprecatedStageManagementRule warns when deprecated stage management is used.
type AwsAPIGatewayDeploymentDeprecatedStageManagementRule struct {
	tflint.DefaultRule

	resourceType           string
	stageNameAttribute     string
	stageDescriptionAttr   string
	canarySettingsBlockTyp string
}

// NewAwsAPIGatewayDeploymentDeprecatedStageManagementRule returns a new rule.
func NewAwsAPIGatewayDeploymentDeprecatedStageManagementRule() *AwsAPIGatewayDeploymentDeprecatedStageManagementRule {
	return &AwsAPIGatewayDeploymentDeprecatedStageManagementRule{
		resourceType:           "aws_api_gateway_deployment",
		stageNameAttribute:     "stage_name",
		stageDescriptionAttr:   "stage_description",
		canarySettingsBlockTyp: "canary_settings",
	}
}

// Name returns the rule name.
func (r *AwsAPIGatewayDeploymentDeprecatedStageManagementRule) Name() string {
	return "awscx_api_gateway_deployment_deprecated_stage_management"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsAPIGatewayDeploymentDeprecatedStageManagementRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsAPIGatewayDeploymentDeprecatedStageManagementRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsAPIGatewayDeploymentDeprecatedStageManagementRule) Link() string {
	return "https://github.com/hashicorp/terraform-provider-aws/issues/39957"
}

// Check warns when deprecated stage management is configured on aws_api_gateway_deployment.
func (r *AwsAPIGatewayDeploymentDeprecatedStageManagementRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.stageNameAttribute},
			{Name: r.stageDescriptionAttr},
		},
		Blocks: []hclext.BlockSchema{
			{Type: r.canarySettingsBlockTyp},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		if attribute, exists := resource.Body.Attributes[r.stageNameAttribute]; exists {
			runner.EmitIssue(
				r,
				"`stage_name` on `aws_api_gateway_deployment` is deprecated; manage stages with `aws_api_gateway_stage` instead.",
				attribute.Expr.Range(),
			)
		}

		if attribute, exists := resource.Body.Attributes[r.stageDescriptionAttr]; exists {
			runner.EmitIssue(
				r,
				"`stage_description` on `aws_api_gateway_deployment` is deprecated; manage stages with `aws_api_gateway_stage` instead.",
				attribute.Expr.Range(),
			)
		}

		for _, block := range resource.Body.Blocks {
			if block.Type != r.canarySettingsBlockTyp {
				continue
			}

			runner.EmitIssue(
				r,
				"`canary_settings` on `aws_api_gateway_deployment` is deprecated; manage stage canary settings with `aws_api_gateway_stage` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
