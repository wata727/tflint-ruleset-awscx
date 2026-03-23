package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLBTargetGroupProtocolVersionNonHTTPRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "tcp target group with protocol version",
			Content: `
resource "aws_lb_target_group" "this" {
  name             = "example"
  port             = 443
  protocol         = "TCP"
  protocol_version = "HTTP2"
  vpc_id           = aws_vpc.main.id
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBTargetGroupProtocolVersionNonHTTPRule(),
					Message: "`protocol_version` can only be set when `protocol` is `HTTP` or `HTTPS`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 22},
						End:      hcl.Pos{Line: 6, Column: 29},
					},
				},
			},
		},
		{
			Name: "lambda target group with protocol version",
			Content: `
resource "aws_lb_target_group" "this" {
  name             = "example"
  target_type      = "lambda"
  protocol_version = "HTTP2"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBTargetGroupProtocolVersionNonHTTPRule(),
					Message: "`protocol_version` can only be set when `protocol` is `HTTP` or `HTTPS`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 22},
						End:      hcl.Pos{Line: 5, Column: 29},
					},
				},
			},
		},
		{
			Name: "https target group with protocol version",
			Content: `
resource "aws_lb_target_group" "this" {
  name             = "example"
  port             = 443
  protocol         = "HTTPS"
  protocol_version = "GRPC"
  vpc_id           = aws_vpc.main.id
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "tcp target group without protocol version",
			Content: `
resource "aws_lb_target_group" "this" {
  name     = "example"
  port     = 443
  protocol = "TCP"
  vpc_id   = aws_vpc.main.id
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown protocol is skipped",
			Content: `
variable "protocol" {
  type = string
}

resource "aws_lb_target_group" "this" {
  name             = "example"
  port             = 443
  protocol         = var.protocol
  protocol_version = "HTTP2"
  vpc_id           = aws_vpc.main.id
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLBTargetGroupProtocolVersionNonHTTPRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
