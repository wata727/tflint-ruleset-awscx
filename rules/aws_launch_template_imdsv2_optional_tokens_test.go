package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLaunchTemplateIMDSv2OptionalTokensRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "optional tokens",
			Content: `
resource "aws_launch_template" "this" {
  name_prefix = "example"

  metadata_options {
    http_tokens = "optional"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLaunchTemplateIMDSv2OptionalTokensRule(),
					Message: "`metadata_options.http_tokens = \"optional\"` allows IMDSv1; prefer `required` to enforce IMDSv2.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 19},
						End:      hcl.Pos{Line: 6, Column: 29},
					},
				},
			},
		},
		{
			Name: "required tokens",
			Content: `
resource "aws_launch_template" "this" {
  name_prefix = "example"

  metadata_options {
    http_tokens = "required"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "metadata_options without http_tokens",
			Content: `
resource "aws_launch_template" "this" {
  name_prefix = "example"

  metadata_options {
    http_endpoint = "enabled"
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLaunchTemplateIMDSv2OptionalTokensRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
