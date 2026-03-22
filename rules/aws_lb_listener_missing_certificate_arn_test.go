package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLBListenerMissingCertificateARNRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "https listener missing certificate",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = 443
  protocol          = "HTTPS"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBListenerMissingCertificateARNRule(),
					Message: "`certificate_arn` must be set when `protocol` is `HTTPS` or `TLS`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 30},
					},
				},
			},
		},
		{
			Name: "tls listener missing certificate",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = 443
  protocol          = "TLS"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBListenerMissingCertificateARNRule(),
					Message: "`certificate_arn` must be set when `protocol` is `HTTPS` or `TLS`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 28},
					},
				},
			},
		},
		{
			Name: "https listener with certificate",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = 443
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate.this.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "http listener without certificate",
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
	}

	rule := NewAwsLBListenerMissingCertificateARNRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
