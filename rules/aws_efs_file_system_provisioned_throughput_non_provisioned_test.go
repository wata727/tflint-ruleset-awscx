package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEFSFileSystemProvisionedThroughputNonProvisionedRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "throughput without throughput mode",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token                  = "example"
  provisioned_throughput_in_mibps = 128
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEFSFileSystemProvisionedThroughputNonProvisionedRule(),
					Message: "`provisioned_throughput_in_mibps` can only be set when `throughput_mode = \"provisioned\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 37},
						End:      hcl.Pos{Line: 4, Column: 40},
					},
				},
			},
		},
		{
			Name: "throughput with elastic mode",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token                  = "example"
  throughput_mode                 = "elastic"
  provisioned_throughput_in_mibps = 128
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEFSFileSystemProvisionedThroughputNonProvisionedRule(),
					Message: "`provisioned_throughput_in_mibps` can only be set when `throughput_mode = \"provisioned\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 37},
						End:      hcl.Pos{Line: 5, Column: 40},
					},
				},
			},
		},
		{
			Name: "throughput with bursting mode",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token                  = "example"
  throughput_mode                 = "bursting"
  provisioned_throughput_in_mibps = 128
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEFSFileSystemProvisionedThroughputNonProvisionedRule(),
					Message: "`provisioned_throughput_in_mibps` can only be set when `throughput_mode = \"provisioned\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 37},
						End:      hcl.Pos{Line: 5, Column: 40},
					},
				},
			},
		},
		{
			Name: "throughput with provisioned mode",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token                  = "example"
  throughput_mode                 = "provisioned"
  provisioned_throughput_in_mibps = 128
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no throughput with elastic mode",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token  = "example"
  throughput_mode = "elastic"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEFSFileSystemProvisionedThroughputNonProvisionedRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
