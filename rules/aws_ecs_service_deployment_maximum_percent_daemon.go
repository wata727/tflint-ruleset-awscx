package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsECSServiceDeploymentMaximumPercentDaemonRule checks invalid deployment_maximum_percent usage on daemon services.
type AwsECSServiceDeploymentMaximumPercentDaemonRule struct {
	tflint.DefaultRule

	resourceType             string
	schedulingStrategyName   string
	deploymentMaximumPercent string
}

// NewAwsECSServiceDeploymentMaximumPercentDaemonRule returns a new rule.
func NewAwsECSServiceDeploymentMaximumPercentDaemonRule() *AwsECSServiceDeploymentMaximumPercentDaemonRule {
	return &AwsECSServiceDeploymentMaximumPercentDaemonRule{
		resourceType:             "aws_ecs_service",
		schedulingStrategyName:   "scheduling_strategy",
		deploymentMaximumPercent: "deployment_maximum_percent",
	}
}

// Name returns the rule name.
func (r *AwsECSServiceDeploymentMaximumPercentDaemonRule) Name() string {
	return "awscx_ecs_service_deployment_maximum_percent_daemon"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsECSServiceDeploymentMaximumPercentDaemonRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsECSServiceDeploymentMaximumPercentDaemonRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsECSServiceDeploymentMaximumPercentDaemonRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service"
}

// Check reports deployment_maximum_percent on daemon-scheduled ECS services.
func (r *AwsECSServiceDeploymentMaximumPercentDaemonRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.schedulingStrategyName},
			{Name: r.deploymentMaximumPercent},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		schedulingStrategy, hasSchedulingStrategy := resource.Body.Attributes[r.schedulingStrategyName]
		if !hasSchedulingStrategy {
			continue
		}

		deploymentMaximumPercent, hasDeploymentMaximumPercent := resource.Body.Attributes[r.deploymentMaximumPercent]
		if !hasDeploymentMaximumPercent {
			continue
		}

		err := runner.EvaluateExpr(schedulingStrategy.Expr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			if !strings.EqualFold(strings.TrimSpace(value.AsString()), "DAEMON") {
				return nil
			}

			runner.EmitIssue(
				r,
				"`deployment_maximum_percent` is not valid when `scheduling_strategy = \"DAEMON\"` on `aws_ecs_service`.",
				deploymentMaximumPercent.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
