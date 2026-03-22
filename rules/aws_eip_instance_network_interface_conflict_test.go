package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEIPInstanceNetworkInterfaceConflictRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "instance and network interface both set",
			Content: `
resource "aws_eip" "this" {
  domain            = "vpc"
  instance          = aws_instance.example.id
  network_interface = aws_network_interface.example.id
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEIPInstanceNetworkInterfaceConflictRule(),
					Message: "`instance` and `network_interface` on `aws_eip` are mutually exclusive; set only one association target.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 55},
					},
				},
			},
		},
		{
			Name: "instance only",
			Content: `
resource "aws_eip" "this" {
  domain   = "vpc"
  instance = aws_instance.example.id
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "network interface only",
			Content: `
resource "aws_eip" "this" {
  domain            = "vpc"
  network_interface = aws_network_interface.example.id
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEIPInstanceNetworkInterfaceConflictRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
