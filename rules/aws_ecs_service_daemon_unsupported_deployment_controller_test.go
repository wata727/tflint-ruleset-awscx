package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsECSServiceDaemonUnsupportedDeploymentControllerRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "daemon scheduling with code deploy controller",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  task_definition     = aws_ecs_task_definition.this.arn
  scheduling_strategy = "DAEMON"

  deployment_controller {
    type = "CODE_DEPLOY"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsECSServiceDaemonUnsupportedDeploymentControllerRule(),
					Message: "`scheduling_strategy = \"DAEMON\"` is not supported when `deployment_controller.type` is `\"CODE_DEPLOY\"` or `\"EXTERNAL\"` on `aws_ecs_service`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 25},
						End:      hcl.Pos{Line: 6, Column: 33},
					},
				},
			},
		},
		{
			Name: "daemon scheduling with external controller",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  scheduling_strategy = "DAEMON"

  deployment_controller {
    type = "EXTERNAL"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsECSServiceDaemonUnsupportedDeploymentControllerRule(),
					Message: "`scheduling_strategy = \"DAEMON\"` is not supported when `deployment_controller.type` is `\"CODE_DEPLOY\"` or `\"EXTERNAL\"` on `aws_ecs_service`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 25},
						End:      hcl.Pos{Line: 5, Column: 33},
					},
				},
			},
		},
		{
			Name: "daemon scheduling with ecs controller",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  scheduling_strategy = "DAEMON"

  deployment_controller {
    type = "ECS"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "replica scheduling with code deploy controller",
			Content: `
resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  scheduling_strategy = "REPLICA"

  deployment_controller {
    type = "CODE_DEPLOY"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown deployment controller type is skipped",
			Content: `
variable "deployment_controller_type" {
  type = string
}

resource "aws_ecs_service" "this" {
  name                = "example"
  cluster             = aws_ecs_cluster.this.id
  scheduling_strategy = "DAEMON"

  deployment_controller {
    type = var.deployment_controller_type
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsECSServiceDaemonUnsupportedDeploymentControllerRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
