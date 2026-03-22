package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEFSFileSystemMissingProvisionedThroughputRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "provisioned mode without throughput value",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token = "example"
  throughput_mode = "provisioned"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEFSFileSystemMissingProvisionedThroughputRule(),
					Message: "`provisioned_throughput_in_mibps` must be set when `throughput_mode = \"provisioned\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 21},
						End:      hcl.Pos{Line: 4, Column: 34},
					},
				},
			},
		},
		{
			Name: "provisioned mode with throughput value",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token                 = "example"
  throughput_mode                = "provisioned"
  provisioned_throughput_in_mibps = 128
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "elastic mode without throughput value",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token = "example"
  throughput_mode = "elastic"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEFSFileSystemMissingProvisionedThroughputRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
