package rules

import (
	"math/big"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule checks Database Insights advanced mode requirements.
type AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule struct {
	tflint.DefaultRule

	resourceType             string
	modeAttribute            string
	enabledAttribute         string
	retentionPeriodAttribute string
	minRetentionPeriod       int
}

// NewAwsDBInstanceDatabaseInsightsAdvancedRequirementsRule returns a new rule.
func NewAwsDBInstanceDatabaseInsightsAdvancedRequirementsRule() *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule {
	return &AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule{
		resourceType:             "aws_db_instance",
		modeAttribute:            "database_insights_mode",
		enabledAttribute:         "performance_insights_enabled",
		retentionPeriodAttribute: "performance_insights_retention_period",
		minRetentionPeriod:       465,
	}
}

// Name returns the rule name.
func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) Name() string {
	return "awscx_db_instance_database_insights_advanced_requirements"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) Link() string {
	return "https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_DatabaseInsights.TurningOnAdvanced.html"
}

// Check reports advanced Database Insights configurations that miss required Performance Insights settings.
func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.modeAttribute},
			{Name: r.enabledAttribute},
			{Name: r.retentionPeriodAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		modeAttr, hasMode := resource.Body.Attributes[r.modeAttribute]
		if !hasMode {
			continue
		}

		mode, ok, err := r.evaluateString(runner, modeAttr.Expr)
		if err != nil {
			return err
		}
		if !ok || !strings.EqualFold(mode, "advanced") {
			continue
		}

		enabledAttr, hasEnabled := resource.Body.Attributes[r.enabledAttribute]
		if !hasEnabled {
			runner.EmitIssue(
				r,
				"`database_insights_mode = \"advanced\"` requires `performance_insights_enabled = true` on `aws_db_instance`.",
				modeAttr.Expr.Range(),
			)
		} else {
			enabled, ok, err := r.evaluateBool(runner, enabledAttr.Expr)
			if err != nil {
				return err
			}
			if ok && !enabled {
				runner.EmitIssue(
					r,
					"`database_insights_mode = \"advanced\"` requires `performance_insights_enabled = true` on `aws_db_instance`.",
					enabledAttr.Expr.Range(),
				)
			}
		}

		retentionAttr, hasRetention := resource.Body.Attributes[r.retentionPeriodAttribute]
		if !hasRetention {
			runner.EmitIssue(
				r,
				"`database_insights_mode = \"advanced\"` requires `performance_insights_retention_period` to be at least `465` on `aws_db_instance`.",
				modeAttr.Expr.Range(),
			)
			continue
		}

		retentionPeriod, ok, err := r.evaluateNumber(runner, retentionAttr.Expr)
		if err != nil {
			return err
		}
		if ok && retentionPeriod < r.minRetentionPeriod {
			runner.EmitIssue(
				r,
				"`database_insights_mode = \"advanced\"` requires `performance_insights_retention_period` to be at least `465` on `aws_db_instance`.",
				retentionAttr.Expr.Range(),
			)
		}
	}

	return nil
}

func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) evaluateString(runner tflint.Runner, expr hcl.Expression) (string, bool, error) {
	value := ""
	resolved := false

	err := runner.EvaluateExpr(expr, func(result cty.Value) error {
		if !result.IsKnown() || result.IsNull() || !result.Type().Equals(cty.String) {
			return nil
		}

		value = result.AsString()
		resolved = true
		return nil
	}, nil)
	if err != nil {
		return "", false, err
	}

	return value, resolved, nil
}

func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) evaluateBool(runner tflint.Runner, expr hcl.Expression) (bool, bool, error) {
	value := false
	resolved := false

	err := runner.EvaluateExpr(expr, func(result cty.Value) error {
		if !result.IsKnown() || result.IsNull() || !result.Type().Equals(cty.Bool) {
			return nil
		}

		value = result.True()
		resolved = true
		return nil
	}, nil)
	if err != nil {
		return false, false, err
	}

	return value, resolved, nil
}

func (r *AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule) evaluateNumber(runner tflint.Runner, expr hcl.Expression) (int, bool, error) {
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
