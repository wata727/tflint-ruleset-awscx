package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsECSServiceDaemonUnsupportedDeploymentControllerRule checks unsupported daemon scheduling with non-ECS deployment controllers.
type AwsECSServiceDaemonUnsupportedDeploymentControllerRule struct {
	tflint.DefaultRule

	resourceType              string
	schedulingStrategyName    string
	deploymentControllerBlock string
	deploymentControllerType  string
}

// NewAwsECSServiceDaemonUnsupportedDeploymentControllerRule returns a new rule.
func NewAwsECSServiceDaemonUnsupportedDeploymentControllerRule() *AwsECSServiceDaemonUnsupportedDeploymentControllerRule {
	return &AwsECSServiceDaemonUnsupportedDeploymentControllerRule{
		resourceType:              "aws_ecs_service",
		schedulingStrategyName:    "scheduling_strategy",
		deploymentControllerBlock: "deployment_controller",
		deploymentControllerType:  "type",
	}
}

// Name returns the rule name.
func (r *AwsECSServiceDaemonUnsupportedDeploymentControllerRule) Name() string {
	return "awscx_ecs_service_daemon_unsupported_deployment_controller"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsECSServiceDaemonUnsupportedDeploymentControllerRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsECSServiceDaemonUnsupportedDeploymentControllerRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsECSServiceDaemonUnsupportedDeploymentControllerRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service"
}

// Check reports DAEMON scheduling on services that explicitly use unsupported deployment controller types.
func (r *AwsECSServiceDaemonUnsupportedDeploymentControllerRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.schedulingStrategyName},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: r.deploymentControllerBlock,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.deploymentControllerType},
					},
				},
			},
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

		isDaemon := false
		err := runner.EvaluateExpr(schedulingStrategy.Expr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			isDaemon = strings.EqualFold(strings.TrimSpace(value.AsString()), "DAEMON")
			return nil
		}, nil)
		if err != nil {
			return err
		}
		if !isDaemon {
			continue
		}

		for _, block := range resource.Body.Blocks {
			controllerType, hasControllerType := block.Body.Attributes[r.deploymentControllerType]
			if !hasControllerType {
				continue
			}

			err := runner.EvaluateExpr(controllerType.Expr, func(value cty.Value) error {
				if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
					return nil
				}

				controller := strings.ToUpper(strings.TrimSpace(value.AsString()))
				if controller != "CODE_DEPLOY" && controller != "EXTERNAL" {
					return nil
				}

				runner.EmitIssue(
					r,
					"`scheduling_strategy = \"DAEMON\"` is not supported when `deployment_controller.type` is `\"CODE_DEPLOY\"` or `\"EXTERNAL\"` on `aws_ecs_service`.",
					schedulingStrategy.Expr.Range(),
				)
				return nil
			}, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
