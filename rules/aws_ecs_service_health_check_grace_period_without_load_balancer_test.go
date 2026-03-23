package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "grace period without load balancer",
			Content: `
resource "aws_ecs_service" "this" {
  name                              = "example"
  cluster                           = aws_ecs_cluster.this.id
  task_definition                   = aws_ecs_task_definition.this.arn
  desired_count                     = 1
  health_check_grace_period_seconds = 60
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule(),
					Message: "`health_check_grace_period_seconds` is only valid when at least one `load_balancer` block is configured on `aws_ecs_service`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 7, Column: 39},
						End:      hcl.Pos{Line: 7, Column: 41},
					},
				},
			},
		},
		{
			Name: "grace period with load balancer",
			Content: `
resource "aws_ecs_service" "this" {
  name                              = "example"
  cluster                           = aws_ecs_cluster.this.id
  task_definition                   = aws_ecs_task_definition.this.arn
  desired_count                     = 1
  health_check_grace_period_seconds = 60

  load_balancer {
    target_group_arn = aws_lb_target_group.this.arn
    container_name   = "app"
    container_port   = 8080
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "service without grace period",
			Content: `
resource "aws_ecs_service" "this" {
  name            = "example"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.this.arn
  desired_count   = 1
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
