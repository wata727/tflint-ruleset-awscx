package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLBListenerMutualAuthenticationVerifyRequirementsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "verify mode without trust store arn",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/app/example/50dc6c495c0c9188"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/example/6d0ecf831eec9f09"
  }

  mutual_authentication {
    mode = "verify"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBListenerMutualAuthenticationVerifyRequirementsRule(),
					Message: "`mutual_authentication.trust_store_arn` must be set when `mutual_authentication.mode` is `\"verify\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 12},
						End:      hcl.Pos{Line: 11, Column: 20},
					},
				},
			},
		},
		{
			Name: "passthrough mode with verify-only attributes",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/app/example/50dc6c495c0c9188"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/example/6d0ecf831eec9f09"
  }

  mutual_authentication {
    mode                             = "passthrough"
    trust_store_arn                  = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:truststore/example/1234567890abcdef"
    ignore_client_certificate_expiry = true
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBListenerMutualAuthenticationVerifyRequirementsRule(),
					Message: "`mutual_authentication.trust_store_arn` is only valid when `mutual_authentication.mode` is `\"verify\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 40},
						End:      hcl.Pos{Line: 12, Column: 134},
					},
				},
				{
					Rule:    NewAwsLBListenerMutualAuthenticationVerifyRequirementsRule(),
					Message: "`mutual_authentication.ignore_client_certificate_expiry` is only valid when `mutual_authentication.mode` is `\"verify\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 13, Column: 40},
						End:      hcl.Pos{Line: 13, Column: 44},
					},
				},
			},
		},
		{
			Name: "off mode with advertised trust store ca names",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/app/example/50dc6c495c0c9188"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/example/6d0ecf831eec9f09"
  }

  mutual_authentication {
    mode                           = "off"
    advertise_trust_store_ca_names = "on"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLBListenerMutualAuthenticationVerifyRequirementsRule(),
					Message: "`mutual_authentication.advertise_trust_store_ca_names` is only valid when `mutual_authentication.mode` is `\"verify\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 38},
						End:      hcl.Pos{Line: 12, Column: 42},
					},
				},
			},
		},
		{
			Name: "verify mode with trust store arn",
			Content: `
resource "aws_lb_listener" "this" {
  load_balancer_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:loadbalancer/app/example/50dc6c495c0c9188"

  default_action {
    type             = "forward"
    target_group_arn = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:targetgroup/example/6d0ecf831eec9f09"
  }

  mutual_authentication {
    mode                             = "verify"
    trust_store_arn                  = "arn:aws:elasticloadbalancing:ap-northeast-1:123456789012:truststore/example/1234567890abcdef"
    advertise_trust_store_ca_names   = "on"
    ignore_client_certificate_expiry = false
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLBListenerMutualAuthenticationVerifyRequirementsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
