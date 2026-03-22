package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEBSVolumeThroughputNonGP3Rule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "throughput without volume type",
			Content: `
resource "aws_ebs_volume" "this" {
  availability_zone = "us-west-2a"
  size              = 40
  throughput        = 250
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEBSVolumeThroughputNonGP3Rule(),
					Message: "`throughput` can only be set when `type = \"gp3\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 26},
					},
				},
			},
		},
		{
			Name: "throughput with gp2",
			Content: `
resource "aws_ebs_volume" "this" {
  availability_zone = "us-west-2a"
  size              = 40
  type              = "gp2"
  throughput        = 250
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEBSVolumeThroughputNonGP3Rule(),
					Message: "`throughput` can only be set when `type = \"gp3\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 23},
						End:      hcl.Pos{Line: 6, Column: 26},
					},
				},
			},
		},
		{
			Name: "throughput with gp3",
			Content: `
resource "aws_ebs_volume" "this" {
  availability_zone = "us-west-2a"
  size              = 40
  type              = "gp3"
  iops              = 3000
  throughput        = 250
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEBSVolumeThroughputNonGP3Rule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
