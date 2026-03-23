package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCloudWatchLogGroupDeliveryRetentionInDaysRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "delivery log group with retention in days",
			Content: `
resource "aws_cloudwatch_log_group" "this" {
  name              = "delivery-log-group"
  log_group_class   = "DELIVERY"
  retention_in_days = 30
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudWatchLogGroupDeliveryRetentionInDaysRule(),
					Message: "`retention_in_days` is ignored when `log_group_class = \"DELIVERY\"`; CloudWatch Logs forces retention to 2 days.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 25},
					},
				},
			},
		},
		{
			Name: "standard log group with retention in days",
			Content: `
resource "aws_cloudwatch_log_group" "this" {
  name              = "standard-log-group"
  log_group_class   = "STANDARD"
  retention_in_days = 30
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "delivery log group without retention in days",
			Content: `
resource "aws_cloudwatch_log_group" "this" {
  name            = "delivery-log-group"
  log_group_class = "DELIVERY"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "delivery log group class with mixed case",
			Content: `
resource "aws_cloudwatch_log_group" "this" {
  name              = "delivery-log-group"
  log_group_class   = "delivery"
  retention_in_days = 7
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudWatchLogGroupDeliveryRetentionInDaysRule(),
					Message: "`retention_in_days` is ignored when `log_group_class = \"DELIVERY\"`; CloudWatch Logs forces retention to 2 days.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 24},
					},
				},
			},
		},
	}

	rule := NewAwsCloudWatchLogGroupDeliveryRetentionInDaysRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
