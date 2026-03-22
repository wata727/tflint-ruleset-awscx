package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSNSTopicFIFONameSuffixRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "fifo topic name missing suffix",
			Content: `
resource "aws_sns_topic" "this" {
  name       = "events"
  fifo_topic = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSNSTopicFIFONameSuffixRule(),
					Message: "`name` must end with `.fifo` when `fifo_topic = true`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 16},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
		},
		{
			Name: "fifo topic name with suffix",
			Content: `
resource "aws_sns_topic" "this" {
  name       = "events.fifo"
  fifo_topic = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "standard topic name without suffix",
			Content: `
resource "aws_sns_topic" "this" {
  name       = "events"
  fifo_topic = false
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "fifo topic with name_prefix only",
			Content: `
resource "aws_sns_topic" "this" {
  name_prefix = "events-"
  fifo_topic  = true
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSNSTopicFIFONameSuffixRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
