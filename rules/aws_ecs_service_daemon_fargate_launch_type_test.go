package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsECSServiceDaemonFargateLaunchTypeRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "daemon scheduling with fargate launch type",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  task_definition     = aws_ecs_task_definition.this.arn
  launch_type         = "FARGATE"
  scheduling_strategy = "DAEMON"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsECSServiceDaemonFargateLaunchTypeRule(),
					Message: "`scheduling_strategy = \"DAEMON\"` is not supported when `launch_type = \"FARGATE\"` on `aws_ecs_service`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 25},
						End:      hcl.Pos{Line: 7, Column: 33},
					},
				},
			},
		},
		{
			Name: "daemon scheduling with ec2 launch type",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  task_definition     = aws_ecs_task_definition.this.arn
  launch_type         = "EC2"
  scheduling_strategy = "DAEMON"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "replica scheduling with fargate launch type",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  task_definition     = aws_ecs_task_definition.this.arn
  launch_type         = "FARGATE"
  scheduling_strategy = "REPLICA"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown launch type is skipped",
			Content: `
variable "launch_type" {
  type = string
}

resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  task_definition     = aws_ecs_task_definition.this.arn
  launch_type         = var.launch_type
  scheduling_strategy = "DAEMON"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsECSServiceDaemonFargateLaunchTypeRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
