package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsECSServiceDeploymentMaximumPercentDaemonRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "deployment maximum percent on daemon service",
			Content: `
resource "aws_ecs_service" "this" {
  name                       = "example"
  cluster                    = aws_ecs_cluster.this.id
  task_definition            = aws_ecs_task_definition.this.arn
  scheduling_strategy        = "DAEMON"
  deployment_maximum_percent = 150
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsECSServiceDeploymentMaximumPercentDaemonRule(),
					Message: "`deployment_maximum_percent` is not valid when `scheduling_strategy = \"DAEMON\"` on `aws_ecs_service`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 32},
						End:      hcl.Pos{Line: 7, Column: 35},
					},
				},
			},
		},
		{
			Name: "deployment maximum percent on replica service",
			Content: `
resource "aws_ecs_service" "this" {
  name                       = "example"
  cluster                    = aws_ecs_cluster.this.id
  task_definition            = aws_ecs_task_definition.this.arn
  scheduling_strategy        = "REPLICA"
  deployment_maximum_percent = 150
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "daemon service without deployment maximum percent",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  task_definition     = aws_ecs_task_definition.this.arn
  scheduling_strategy = "DAEMON"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown scheduling strategy is skipped",
			Content: `
variable "scheduling_strategy" {
  type = string
}

resource "aws_ecs_service" "this" {
  name                       = "example"
  cluster                    = aws_ecs_cluster.this.id
  task_definition            = aws_ecs_task_definition.this.arn
  scheduling_strategy        = var.scheduling_strategy
  deployment_maximum_percent = 150
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsECSServiceDeploymentMaximumPercentDaemonRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
