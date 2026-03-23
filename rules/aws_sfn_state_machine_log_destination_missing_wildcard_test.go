package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSFNStateMachineLogDestinationMissingWildcardRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "log destination without wildcard suffix",
			Content: `
resource "aws_sfn_state_machine" "this" {
  name     = "example"
  role_arn = aws_iam_role.this.arn
  definition = jsonencode({
    StartAt = "Pass"
    States = {
      Pass = {
        Type = "Pass"
        End  = true
      }
    }
  })

  logging_configuration {
    level           = "ERROR"
    include_execution_data = false
    log_destination = "arn:aws:logs:ap-northeast-1:123456789012:log-group:/aws/vendedlogs/states/example"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsSFNStateMachineLogDestinationMissingWildcardRule(),
					Message: "`logging_configuration.log_destination` must end with `:*` for Step Functions CloudWatch Logs destinations.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 18, Column: 23},
						End:      hcl.Pos{Line: 18, Column: 106},
					},
				},
			},
		},
		{
			Name: "log destination with wildcard suffix",
			Content: `
resource "aws_sfn_state_machine" "this" {
  name     = "example"
  role_arn = aws_iam_role.this.arn
  definition = jsonencode({
    StartAt = "Pass"
    States = {
      Pass = {
        Type = "Pass"
        End  = true
      }
    }
  })

  logging_configuration {
    level           = "ERROR"
    include_execution_data = false
    log_destination = "arn:aws:logs:ap-northeast-1:123456789012:log-group:/aws/vendedlogs/states/example:*"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown log destination is skipped",
			Content: `
variable "log_destination" {
  type = string
}

resource "aws_sfn_state_machine" "this" {
  name     = "example"
  role_arn = aws_iam_role.this.arn
  definition = jsonencode({
    StartAt = "Pass"
    States = {
      Pass = {
        Type = "Pass"
        End  = true
      }
    }
  })

  logging_configuration {
    level                  = "ERROR"
    include_execution_data = false
    log_destination        = var.log_destination
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "state machine without logging configuration",
			Content: `
resource "aws_sfn_state_machine" "this" {
  name     = "example"
  role_arn = aws_iam_role.this.arn
  definition = jsonencode({
    StartAt = "Pass"
    States = {
      Pass = {
        Type = "Pass"
        End  = true
      }
    }
  })
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsSFNStateMachineLogDestinationMissingWildcardRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
