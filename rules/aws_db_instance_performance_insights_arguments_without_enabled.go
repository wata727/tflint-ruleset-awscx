package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule checks invalid Performance Insights argument combinations.
type AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule struct {
	tflint.DefaultRule

	resourceType             string
	enabledAttribute         string
	kmsKeyAttribute          string
	retentionPeriodAttribute string
}

// NewAwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule returns a new rule.
func NewAwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule() *AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule {
	return &AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule{
		resourceType:             "aws_db_instance",
		enabledAttribute:         "performance_insights_enabled",
		kmsKeyAttribute:          "performance_insights_kms_key_id",
		retentionPeriodAttribute: "performance_insights_retention_period",
	}
}

// Name returns the rule name.
func (r *AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule) Name() string {
	return "awscx_db_instance_performance_insights_arguments_without_enabled"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance"
}

// Check reports Performance Insights arguments configured without enabling Performance Insights.
func (r *AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.enabledAttribute},
			{Name: r.kmsKeyAttribute},
			{Name: r.retentionPeriodAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		enabled, hasEnabled := resource.Body.Attributes[r.enabledAttribute]
		kmsKey, hasKMSKey := resource.Body.Attributes[r.kmsKeyAttribute]
		retentionPeriod, hasRetentionPeriod := resource.Body.Attributes[r.retentionPeriodAttribute]

		if !hasKMSKey && !hasRetentionPeriod {
			continue
		}

		if !hasEnabled {
			r.emitIssues(runner, hasKMSKey, kmsKey, hasRetentionPeriod, retentionPeriod)
			continue
		}

		err := runner.EvaluateExpr(enabled.Expr, func(value bool) error {
			if value {
				return nil
			}

			r.emitIssues(runner, hasKMSKey, kmsKey, hasRetentionPeriod, retentionPeriod)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule) emitIssues(
	runner tflint.Runner,
	hasKMSKey bool,
	kmsKey *hclext.Attribute,
	hasRetentionPeriod bool,
	retentionPeriod *hclext.Attribute,
) {
	if hasKMSKey {
		runner.EmitIssue(
			r,
			"`performance_insights_kms_key_id` cannot be set unless `performance_insights_enabled = true` on `aws_db_instance`.",
			kmsKey.Expr.Range(),
		)
	}

	if hasRetentionPeriod {
		runner.EmitIssue(
			r,
			"`performance_insights_retention_period` cannot be set unless `performance_insights_enabled = true` on `aws_db_instance`.",
			retentionPeriod.Expr.Range(),
		)
	}
}
