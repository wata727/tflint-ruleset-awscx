package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSQSQueuePolicyMissingCurrentVersionRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "missing version in heredoc policy",
			Content: `
resource "aws_sqs_queue_policy" "this" {
  queue_url = aws_sqs_queue.this.id
  policy = <<POLICY
{
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "*"
    }
  ]
}
POLICY
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSQSQueuePolicyMissingCurrentVersionRule(),
					Message: "`policy` on `aws_sqs_queue_policy` must set top-level `Version` to `2012-10-17`; AWS may time out without it.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 12},
						End:      hcl.Pos{Line: 15, Column: 7},
					},
				},
			},
		},
		{
			Name: "legacy version in jsonencode policy",
			Content: `
resource "aws_sqs_queue_policy" "this" {
  queue_url = aws_sqs_queue.this.id
  policy = jsonencode({
    Version = "2008-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = "*"
      Action = "sqs:SendMessage"
      Resource = "*"
    }]
  })
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSQSQueuePolicyMissingCurrentVersionRule(),
					Message: "`policy` on `aws_sqs_queue_policy` must set top-level `Version` to `2012-10-17`; AWS may time out without it.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 12},
						End:      hcl.Pos{Line: 12, Column: 5},
					},
				},
			},
		},
		{
			Name: "current version in jsonencode policy",
			Content: `
resource "aws_sqs_queue_policy" "this" {
  queue_url = aws_sqs_queue.this.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = "*"
      Action = "sqs:SendMessage"
      Resource = "*"
    }]
  })
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown policy is skipped",
			Content: `
variable "policy" {
  type = string
}

resource "aws_sqs_queue_policy" "this" {
  queue_url = aws_sqs_queue.this.id
  policy    = var.policy
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSQSQueuePolicyMissingCurrentVersionRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
