package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLBTargetGroupLambdaTopLevelAttributesRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "lambda target group with port protocol and vpc_id",
			Content: `
resource "aws_lb_target_group" "this" {
  name        = "example"
  target_type = "lambda"
  port        = 443
  protocol    = "HTTPS"
  vpc_id      = aws_vpc.main.id
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBTargetGroupLambdaTopLevelAttributesRule(),
					Message: "`port` does not apply when `target_type = \"lambda\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 17},
						End:      hcl.Pos{Line: 5, Column: 20},
					},
				},
				{
					Rule:    NewAwsLBTargetGroupLambdaTopLevelAttributesRule(),
					Message: "`protocol` does not apply when `target_type = \"lambda\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 17},
						End:      hcl.Pos{Line: 6, Column: 24},
					},
				},
				{
					Rule:    NewAwsLBTargetGroupLambdaTopLevelAttributesRule(),
					Message: "`vpc_id` does not apply when `target_type = \"lambda\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 17},
						End:      hcl.Pos{Line: 7, Column: 32},
					},
				},
			},
		},
		{
			Name: "lambda target group without top level attributes",
			Content: `
resource "aws_lb_target_group" "this" {
  name        = "example"
  target_type = "lambda"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "instance target group with port protocol and vpc_id",
			Content: `
resource "aws_lb_target_group" "this" {
  name        = "example"
  target_type = "instance"
  port        = 443
  protocol    = "HTTPS"
  vpc_id      = aws_vpc.main.id
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown target type is skipped",
			Content: `
variable "target_type" {
  type = string
}

resource "aws_lb_target_group" "this" {
  name        = "example"
  target_type = var.target_type
  port        = 443
  protocol    = "HTTPS"
  vpc_id      = aws_vpc.main.id
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLBTargetGroupLambdaTopLevelAttributesRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
