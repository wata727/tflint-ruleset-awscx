package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAPIGatewayDeploymentDeprecatedStageManagementRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "deprecated stage_name and stage_description attributes",
			Content: `
resource "aws_api_gateway_deployment" "this" {
  rest_api_id       = aws_api_gateway_rest_api.this.id
  triggers          = { redeployment = "1" }
  stage_name        = "prod"
  stage_description = "Production stage"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayDeploymentDeprecatedStageManagementRule(),
					Message: "`stage_name` on `aws_api_gateway_deployment` is deprecated; manage stages with `aws_api_gateway_stage` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 29},
					},
				},
				{
					Rule:    NewAwsAPIGatewayDeploymentDeprecatedStageManagementRule(),
					Message: "`stage_description` on `aws_api_gateway_deployment` is deprecated; manage stages with `aws_api_gateway_stage` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 23},
						End:      hcl.Pos{Line: 6, Column: 41},
					},
				},
			},
		},
		{
			Name: "deprecated canary_settings block",
			Content: `
resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id

  canary_settings {
    percent_traffic = 10
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAPIGatewayDeploymentDeprecatedStageManagementRule(),
					Message: "`canary_settings` on `aws_api_gateway_deployment` is deprecated; manage stage canary settings with `aws_api_gateway_stage` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 18},
					},
				},
			},
		},
		{
			Name: "separate stage resource",
			Content: `
resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
}

resource "aws_api_gateway_stage" "this" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  deployment_id = aws_api_gateway_deployment.this.id
  stage_name    = "prod"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAPIGatewayDeploymentDeprecatedStageManagementRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
