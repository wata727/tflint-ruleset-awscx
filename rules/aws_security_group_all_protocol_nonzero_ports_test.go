package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupAllProtocolNonzeroPortsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "ingress all protocol with nonzero port range",
			Content: `
resource "aws_security_group" "this" {
  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "-1"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupAllProtocolNonzeroPortsRule(),
					Message: "`from_port` and `to_port` must both be `0` when `protocol = \"-1\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 17},
						End:      hcl.Pos{Line: 6, Column: 21},
					},
				},
			},
		},
		{
			Name: "egress all protocol with nonzero end port",
			Content: `
resource "aws_security_group" "this" {
  egress {
    from_port = 0
    to_port   = 65535
    protocol  = "-1"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSecurityGroupAllProtocolNonzeroPortsRule(),
					Message: "`from_port` and `to_port` must both be `0` when `protocol = \"-1\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 17},
						End:      hcl.Pos{Line: 6, Column: 21},
					},
				},
			},
		},
		{
			Name: "all protocol with zero ports",
			Content: `
resource "aws_security_group" "this" {
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "specific protocol with nonzero ports",
			Content: `
resource "aws_security_group" "this" {
  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "tcp"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown port expressions are skipped",
			Content: `
variable "port" {
  type = number
}

resource "aws_security_group" "this" {
  ingress {
    from_port = var.port
    to_port   = var.port
    protocol  = "-1"
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSecurityGroupAllProtocolNonzeroPortsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
