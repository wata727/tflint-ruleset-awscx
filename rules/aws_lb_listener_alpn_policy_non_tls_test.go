package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLBListenerALPNPolicyNonTLSRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "https listener with alpn policy",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = 443
  protocol          = "HTTPS"
  alpn_policy       = "HTTP2Preferred"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBListenerALPNPolicyNonTLSRule(),
					Message: "`alpn_policy` can only be set when `protocol` is `TLS`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 23},
						End:      hcl.Pos{Line: 6, Column: 39},
					},
				},
			},
		},
		{
			Name: "tls listener with alpn policy",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = 443
  protocol          = "TLS"
  alpn_policy       = "HTTP2Preferred"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "http listener without alpn policy",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "listener with alpn policy and unknown protocol is skipped",
			Content: `
variable "listener_protocol" {
  type = string
}

resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = 443
  protocol          = var.listener_protocol
  alpn_policy       = "HTTP2Preferred"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLBListenerALPNPolicyNonTLSRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
