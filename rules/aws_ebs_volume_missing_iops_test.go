package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEBSVolumeMissingIOPSRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "io1 without iops",
			Content: `
resource "aws_ebs_volume" "this" {
  availability_zone = "ap-northeast-1a"
  size              = 100
  type              = "io1"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEBSVolumeMissingIOPSRule(),
					Message: "`iops` must be set when `type` is `io1` or `io2`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 28},
					},
				},
			},
		},
		{
			Name: "io2 with iops",
			Content: `
resource "aws_ebs_volume" "this" {
  availability_zone = "ap-northeast-1a"
  size              = 100
  type              = "io2"
  iops              = 3000
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "gp3 without iops",
			Content: `
resource "aws_ebs_volume" "this" {
  availability_zone = "ap-northeast-1a"
  size              = 100
  type              = "gp3"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEBSVolumeMissingIOPSRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
