package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule checks invalid ECS grace period usage.
type AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule struct {
	tflint.DefaultRule

	resourceType       string
	attributeName      string
	loadBalancerBlocks string
}

// NewAwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule returns a new rule.
func NewAwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule() *AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule {
	return &AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule{
		resourceType:       "aws_ecs_service",
		attributeName:      "health_check_grace_period_seconds",
		loadBalancerBlocks: "load_balancer",
	}
}

// Name returns the rule name.
func (r *AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule) Name() string {
	return "awscx_ecs_service_health_check_grace_period_without_load_balancer"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service"
}

// Check reports grace periods configured without any load balancer block.
func (r *AwsECSServiceHealthCheckGracePeriodWithoutLoadBalancerRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.attributeName},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: r.loadBalancerBlocks,
				Body: &hclext.BodySchema{},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, hasAttribute := resource.Body.Attributes[r.attributeName]
		if !hasAttribute {
			continue
		}

		if len(resource.Body.Blocks) > 0 {
			continue
		}

		runner.EmitIssue(
			r,
			"`health_check_grace_period_seconds` is only valid when at least one `load_balancer` block is configured on `aws_ecs_service`.",
			attribute.Expr.Range(),
		)
	}

	return nil
}
