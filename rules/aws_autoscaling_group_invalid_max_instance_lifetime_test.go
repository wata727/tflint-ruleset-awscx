package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "zero lifetime is allowed",
			Content: `
resource "aws_autoscaling_group" "this" {
  name                  = "example"
  max_size              = 2
  min_size              = 1
  desired_capacity      = 1
  vpc_zone_identifier   = [aws_subnet.example.id]
  max_instance_lifetime = 0
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "documented minimum lifetime is allowed",
			Content: `
resource "aws_autoscaling_group" "this" {
  name                  = "example"
  max_size              = 2
  min_size              = 1
  desired_capacity      = 1
  vpc_zone_identifier   = [aws_subnet.example.id]
  max_instance_lifetime = 86400
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "value below minimum is invalid",
			Content: `
resource "aws_autoscaling_group" "this" {
  name                  = "example"
  max_size              = 2
  min_size              = 1
  desired_capacity      = 1
  vpc_zone_identifier   = [aws_subnet.example.id]
  max_instance_lifetime = 86399
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAutoscalingGroupInvalidMaxInstanceLifetimeRule(),
					Message: "`max_instance_lifetime` must be 0 or between 86400 and 31536000 seconds on `aws_autoscaling_group`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 8, Column: 27},
						End:      hcl.Pos{Line: 8, Column: 32},
					},
				},
			},
		},
		{
			Name: "value above maximum is invalid",
			Content: `
resource "aws_autoscaling_group" "this" {
  name                  = "example"
  max_size              = 2
  min_size              = 1
  desired_capacity      = 1
  vpc_zone_identifier   = [aws_subnet.example.id]
  max_instance_lifetime = 31536001
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsAutoscalingGroupInvalidMaxInstanceLifetimeRule(),
					Message: "`max_instance_lifetime` must be 0 or between 86400 and 31536000 seconds on `aws_autoscaling_group`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 8, Column: 27},
						End:      hcl.Pos{Line: 8, Column: 35},
					},
				},
			},
		},
		{
			Name: "unknown value is skipped",
			Content: `
variable "lifetime" {
  type = number
}

resource "aws_autoscaling_group" "this" {
  name                  = "example"
  max_size              = 2
  min_size              = 1
  desired_capacity      = 1
  vpc_zone_identifier   = [aws_subnet.example.id]
  max_instance_lifetime = var.lifetime
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsAutoscalingGroupInvalidMaxInstanceLifetimeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
