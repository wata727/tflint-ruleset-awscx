package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsInstanceDeprecatedNetworkInterfaceRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "deprecated network_interface block",
			Content: `
resource "aws_instance" "this" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"

  network_interface {
    device_index         = 0
    network_interface_id = "eni-12345678"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsInstanceDeprecatedNetworkInterfaceRule(),
					Message: "`network_interface` on `aws_instance` is deprecated; use `primary_network_interface` for the primary ENI and `aws_network_interface_attachment` for additional ENIs instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 20},
					},
				},
			},
		},
		{
			Name: "primary_network_interface block",
			Content: `
resource "aws_instance" "this" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"

  primary_network_interface {
    network_interface_id = "eni-12345678"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "instance without network interface blocks",
			Content: `
resource "aws_instance" "this" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"
  subnet_id     = "subnet-12345678"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsInstanceDeprecatedNetworkInterfaceRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
