package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "archive policy without fifo topic",
			Content: `
resource "aws_sns_topic" "this" {
  name           = "updates"
  archive_policy = jsonencode({ MessageRetentionPeriod = "30" })
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule(),
					Message: "`archive_policy` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 20},
						End:      hcl.Pos{Line: 4, Column: 65},
					},
				},
			},
		},
		{
			Name: "content based deduplication without fifo topic",
			Content: `
resource "aws_sns_topic" "this" {
  name                        = "updates"
  content_based_deduplication = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule(),
					Message: "`content_based_deduplication` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 33},
						End:      hcl.Pos{Line: 4, Column: 37},
					},
				},
			},
		},
		{
			Name: "fifo throughput scope with explicit false fifo topic",
			Content: `
resource "aws_sns_topic" "this" {
  name                  = "updates"
  fifo_topic            = false
  fifo_throughput_scope = "MessageGroup"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule(),
					Message: "`fifo_throughput_scope` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 27},
						End:      hcl.Pos{Line: 5, Column: 41},
					},
				},
			},
		},
		{
			Name: "multiple fifo only attributes without fifo topic",
			Content: `
resource "aws_sns_topic" "this" {
  name                        = "updates"
  archive_policy              = jsonencode({ MessageRetentionPeriod = "30" })
  content_based_deduplication = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule(),
					Message: "`archive_policy` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 33},
						End:      hcl.Pos{Line: 4, Column: 78},
					},
				},
				{
					Rule:    NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule(),
					Message: "`content_based_deduplication` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 33},
						End:      hcl.Pos{Line: 5, Column: 37},
					},
				},
			},
		},
		{
			Name: "fifo topic with fifo only attributes",
			Content: `
resource "aws_sns_topic" "this" {
  name                        = "updates.fifo"
  fifo_topic                  = true
  archive_policy              = jsonencode({ MessageRetentionPeriod = "30" })
  content_based_deduplication = true
  fifo_throughput_scope       = "MessageGroup"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
