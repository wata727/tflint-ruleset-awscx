package rules

import (
	"math/big"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule checks Blue/Green backup retention prerequisites.
type AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule struct {
	tflint.DefaultRule

	resourceType               string
	backupRetentionPeriodAttr  string
	blueGreenUpdateBlockType   string
	blueGreenUpdateEnabledAttr string
}

// NewAwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule returns a new rule.
func NewAwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule() *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule {
	return &AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule{
		resourceType:               "aws_db_instance",
		backupRetentionPeriodAttr:  "backup_retention_period",
		blueGreenUpdateBlockType:   "blue_green_update",
		blueGreenUpdateEnabledAttr: "enabled",
	}
}

// Name returns the rule name.
func (r *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule) Name() string {
	return "awscx_db_instance_blue_green_update_without_backup_retention"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule) Link() string {
	return "https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/blue-green-deployments-creating.html"
}

// Check reports aws_db_instance Blue/Green configurations without automated backups enabled.
func (r *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.backupRetentionPeriodAttr},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blueGreenUpdateBlockType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.blueGreenUpdateEnabledAttr},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		backupRetentionAttr, hasBackupRetention := resource.Body.Attributes[r.backupRetentionPeriodAttr]

		for _, block := range resource.Body.Blocks {
			if block.Type != r.blueGreenUpdateBlockType {
				continue
			}

			enabledAttr, hasEnabled := block.Body.Attributes[r.blueGreenUpdateEnabledAttr]
			if !hasEnabled {
				continue
			}

			enabled, ok, err := r.evaluateBool(runner, enabledAttr.Expr)
			if err != nil {
				return err
			}
			if !ok || !enabled {
				continue
			}

			if !hasBackupRetention {
				runner.EmitIssue(
					r,
					"`blue_green_update.enabled = true` requires `backup_retention_period` greater than `0` on `aws_db_instance`.",
					enabledAttr.Expr.Range(),
				)
				continue
			}

			backupRetentionPeriod, ok, err := r.evaluateInteger(runner, backupRetentionAttr.Expr)
			if err != nil {
				return err
			}
			if ok && backupRetentionPeriod <= 0 {
				runner.EmitIssue(
					r,
					"`blue_green_update.enabled = true` requires `backup_retention_period` greater than `0` on `aws_db_instance`.",
					backupRetentionAttr.Expr.Range(),
				)
			}
		}
	}

	return nil
}

func (r *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule) evaluateBool(runner tflint.Runner, expr hcl.Expression) (bool, bool, error) {
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

func (r *AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule) evaluateInteger(runner tflint.Runner, expr hcl.Expression) (int, bool, error) {
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
