package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLBTargetGroupMatcherNonHTTPHealthCheckRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "tcp health check with matcher",
			Content: `
resource "aws_lb_target_group" "this" {
  name     = "example"
  port     = 80
  protocol = "TCP"
  vpc_id   = aws_vpc.main.id

  health_check {
    protocol = "TCP"
    matcher  = "200-399"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBTargetGroupMatcherNonHTTPHealthCheckRule(),
					Message: "`health_check.matcher` can only be set when `health_check.protocol` is `HTTP` or `HTTPS`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 16},
						End:      hcl.Pos{Line: 10, Column: 25},
					},
				},
			},
		},
		{
			Name: "http health check with matcher",
			Content: `
resource "aws_lb_target_group" "this" {
  name     = "example"
  port     = 80
  protocol = "TCP"
  vpc_id   = aws_vpc.main.id

  health_check {
    protocol = "HTTP"
    matcher  = "200-399"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "tcp health check without matcher",
			Content: `
resource "aws_lb_target_group" "this" {
  name     = "example"
  port     = 80
  protocol = "TCP"
  vpc_id   = aws_vpc.main.id

  health_check {
    protocol = "TCP"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "matcher without health check protocol is skipped",
			Content: `
resource "aws_lb_target_group" "this" {
  name     = "example"
  port     = 80
  protocol = "TCP"
  vpc_id   = aws_vpc.main.id

  health_check {
    matcher = "200-399"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "lambda target group matcher is allowed",
			Content: `
resource "aws_lb_target_group" "this" {
  name        = "example"
  target_type = "lambda"

  health_check {
    matcher = "200-499"
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLBTargetGroupMatcherNonHTTPHealthCheckRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
