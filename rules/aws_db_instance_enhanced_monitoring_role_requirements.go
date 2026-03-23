package rules

import (
	"math/big"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsDBInstanceEnhancedMonitoringRoleRequirementsRule checks Enhanced Monitoring argument dependencies.
type AwsDBInstanceEnhancedMonitoringRoleRequirementsRule struct {
	tflint.DefaultRule

	resourceType           string
	monitoringIntervalAttr string
	monitoringRoleARNAttr  string
}

// NewAwsDBInstanceEnhancedMonitoringRoleRequirementsRule returns a new rule.
func NewAwsDBInstanceEnhancedMonitoringRoleRequirementsRule() *AwsDBInstanceEnhancedMonitoringRoleRequirementsRule {
	return &AwsDBInstanceEnhancedMonitoringRoleRequirementsRule{
		resourceType:           "aws_db_instance",
		monitoringIntervalAttr: "monitoring_interval",
		monitoringRoleARNAttr:  "monitoring_role_arn",
	}
}

// Name returns the rule name.
func (r *AwsDBInstanceEnhancedMonitoringRoleRequirementsRule) Name() string {
	return "awscx_db_instance_enhanced_monitoring_role_requirements"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstanceEnhancedMonitoringRoleRequirementsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstanceEnhancedMonitoringRoleRequirementsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstanceEnhancedMonitoringRoleRequirementsRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance"
}

// Check reports mismatched Enhanced Monitoring interval and role settings.
func (r *AwsDBInstanceEnhancedMonitoringRoleRequirementsRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.monitoringIntervalAttr},
			{Name: r.monitoringRoleARNAttr},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		monitoringIntervalAttr, hasMonitoringInterval := resource.Body.Attributes[r.monitoringIntervalAttr]
		monitoringRoleARNAttr, hasMonitoringRoleARN := resource.Body.Attributes[r.monitoringRoleARNAttr]
		if !hasMonitoringInterval {
			continue
		}

		monitoringInterval, ok, err := r.evaluateInteger(runner, monitoringIntervalAttr.Expr)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		if monitoringInterval == 0 {
			if hasMonitoringRoleARN {
				runner.EmitIssue(
					r,
					"`monitoring_role_arn` cannot be set when `monitoring_interval = 0` on `aws_db_instance`.",
					monitoringRoleARNAttr.Expr.Range(),
				)
			}
			continue
		}

		if hasMonitoringRoleARN {
			continue
		}

		runner.EmitIssue(
			r,
			"`monitoring_role_arn` must be set when `monitoring_interval` enables Enhanced Monitoring on `aws_db_instance`.",
			monitoringIntervalAttr.Expr.Range(),
		)
	}

	return nil
}

func (r *AwsDBInstanceEnhancedMonitoringRoleRequirementsRule) evaluateInteger(runner tflint.Runner, expr hcl.Expression) (int, bool, error) {
	value := 0
	resolved := false

	err := runner.EvaluateExpr(expr, func(result cty.Value) error {
		if !result.IsKnown() || result.IsNull() || !result.Type().Equals(cty.Number) {
			return nil
		}

		number := result.AsBigFloat()
		intValue, accuracy := number.Int64()
		if accuracy != big.Exact {
			return nil
		}

		value = int(intValue)
		resolved = true
		return nil
	}, nil)
	if err != nil {
		return 0, false, err
	}

	return value, resolved, nil
}
