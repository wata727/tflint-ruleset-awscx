package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsCloudWatchLogGroupDeliveryRetentionInDaysRule warns when DELIVERY log groups set retention_in_days.
type AwsCloudWatchLogGroupDeliveryRetentionInDaysRule struct {
	tflint.DefaultRule

	resourceType        string
	logGroupClassAttr   string
	retentionInDaysAttr string
}

// NewAwsCloudWatchLogGroupDeliveryRetentionInDaysRule returns a new rule.
func NewAwsCloudWatchLogGroupDeliveryRetentionInDaysRule() *AwsCloudWatchLogGroupDeliveryRetentionInDaysRule {
	return &AwsCloudWatchLogGroupDeliveryRetentionInDaysRule{
		resourceType:        "aws_cloudwatch_log_group",
		logGroupClassAttr:   "log_group_class",
		retentionInDaysAttr: "retention_in_days",
	}
}

// Name returns the rule name.
func (r *AwsCloudWatchLogGroupDeliveryRetentionInDaysRule) Name() string {
	return "awscx_cloudwatch_log_group_delivery_retention_in_days"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsCloudWatchLogGroupDeliveryRetentionInDaysRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsCloudWatchLogGroupDeliveryRetentionInDaysRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsCloudWatchLogGroupDeliveryRetentionInDaysRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group"
}

// Check warns when DELIVERY log groups explicitly configure retention_in_days.
func (r *AwsCloudWatchLogGroupDeliveryRetentionInDaysRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.logGroupClassAttr},
			{Name: r.retentionInDaysAttr},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		logGroupClass, hasLogGroupClass := resource.Body.Attributes[r.logGroupClassAttr]
		if !hasLogGroupClass {
			continue
		}

		retentionInDays, hasRetentionInDays := resource.Body.Attributes[r.retentionInDaysAttr]
		if !hasRetentionInDays {
			continue
		}

		err := runner.EvaluateExpr(logGroupClass.Expr, func(value string) error {
			if !strings.EqualFold(strings.TrimSpace(value), "DELIVERY") {
				return nil
			}

			runner.EmitIssue(
				r,
				"`retention_in_days` is ignored when `log_group_class = \"DELIVERY\"`; CloudWatch Logs forces retention to 2 days.",
				retentionInDays.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
