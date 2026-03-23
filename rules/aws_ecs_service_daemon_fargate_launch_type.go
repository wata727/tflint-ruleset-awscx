package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsECSServiceDaemonFargateLaunchTypeRule checks unsupported daemon scheduling on Fargate services.
type AwsECSServiceDaemonFargateLaunchTypeRule struct {
	tflint.DefaultRule

	resourceType            string
	launchTypeAttribute     string
	schedulingTypeAttribute string
}

// NewAwsECSServiceDaemonFargateLaunchTypeRule returns a new rule.
func NewAwsECSServiceDaemonFargateLaunchTypeRule() *AwsECSServiceDaemonFargateLaunchTypeRule {
	return &AwsECSServiceDaemonFargateLaunchTypeRule{
		resourceType:            "aws_ecs_service",
		launchTypeAttribute:     "launch_type",
		schedulingTypeAttribute: "scheduling_strategy",
	}
}

// Name returns the rule name.
func (r *AwsECSServiceDaemonFargateLaunchTypeRule) Name() string {
	return "awscx_ecs_service_daemon_fargate_launch_type"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsECSServiceDaemonFargateLaunchTypeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsECSServiceDaemonFargateLaunchTypeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsECSServiceDaemonFargateLaunchTypeRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service"
}

// Check reports DAEMON scheduling on services that explicitly use the FARGATE launch type.
func (r *AwsECSServiceDaemonFargateLaunchTypeRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.launchTypeAttribute},
			{Name: r.schedulingTypeAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		launchType, hasLaunchType := resource.Body.Attributes[r.launchTypeAttribute]
		schedulingStrategy, hasSchedulingStrategy := resource.Body.Attributes[r.schedulingTypeAttribute]
		if !hasLaunchType || !hasSchedulingStrategy {
			continue
		}

		isFargate := false
		err := runner.EvaluateExpr(launchType.Expr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			isFargate = strings.EqualFold(strings.TrimSpace(value.AsString()), "FARGATE")
			return nil
		}, nil)
		if err != nil {
			return err
		}
		if !isFargate {
			continue
		}

		err = runner.EvaluateExpr(schedulingStrategy.Expr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			if !strings.EqualFold(strings.TrimSpace(value.AsString()), "DAEMON") {
				return nil
			}

			runner.EmitIssue(
				r,
				"`scheduling_strategy = \"DAEMON\"` is not supported when `launch_type = \"FARGATE\"` on `aws_ecs_service`.",
				schedulingStrategy.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
