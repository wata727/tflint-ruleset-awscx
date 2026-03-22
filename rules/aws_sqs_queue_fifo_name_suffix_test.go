package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSQSQueueFIFONameSuffixRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "fifo queue name missing suffix",
			Content: `
resource "aws_sqs_queue" "this" {
  name       = "orders"
  fifo_queue = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSQSQueueFIFONameSuffixRule(),
					Message: "`name` must end with `.fifo` when `fifo_queue = true`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 16},
						End:      hcl.Pos{Line: 3, Column: 24},
					},
				},
			},
		},
		{
			Name: "fifo queue name with suffix",
			Content: `
resource "aws_sqs_queue" "this" {
  name       = "orders.fifo"
  fifo_queue = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "standard queue name without suffix",
			Content: `
resource "aws_sqs_queue" "this" {
  name       = "orders"
  fifo_queue = false
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "fifo queue with name_prefix only",
			Content: `
resource "aws_sqs_queue" "this" {
  name_prefix = "orders-"
  fifo_queue  = true
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSQSQueueFIFONameSuffixRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
