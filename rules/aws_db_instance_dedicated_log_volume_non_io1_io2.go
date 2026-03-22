package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule checks that dedicated_log_volume is only used with io1 or io2 storage.
type AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule struct {
	tflint.DefaultRule

	resourceType                string
	dedicatedLogVolumeAttribute string
	storageTypeAttribute        string
}

// NewAwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule returns a new rule.
func NewAwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule() *AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule {
	return &AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule{
		resourceType:                "aws_db_instance",
		dedicatedLogVolumeAttribute: "dedicated_log_volume",
		storageTypeAttribute:        "storage_type",
	}
}

// Name returns the rule name.
func (r *AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule) Name() string {
	return "awscx_db_instance_dedicated_log_volume_non_io1_io2"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance"
}

// Check reports when dedicated_log_volume is configured without io1 or io2 storage.
func (r *AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.dedicatedLogVolumeAttribute},
			{Name: r.storageTypeAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		dedicatedLogVolume, hasDedicatedLogVolume := resource.Body.Attributes[r.dedicatedLogVolumeAttribute]
		if !hasDedicatedLogVolume {
			continue
		}

		storageType, hasStorageType := resource.Body.Attributes[r.storageTypeAttribute]

		err := runner.EvaluateExpr(dedicatedLogVolume.Expr, func(enabled bool) error {
			if !enabled {
				return nil
			}

			if !hasStorageType {
				runner.EmitIssue(
					r,
					"`dedicated_log_volume` can only be set when `storage_type` is `io1` or `io2`.",
					dedicatedLogVolume.Expr.Range(),
				)
				return nil
			}

			return runner.EvaluateExpr(storageType.Expr, func(value string) error {
				switch strings.ToLower(value) {
				case "io1", "io2":
					return nil
				default:
					runner.EmitIssue(
						r,
						"`dedicated_log_volume` can only be set when `storage_type` is `io1` or `io2`.",
						dedicatedLogVolume.Expr.Range(),
					)
					return nil
				}
			}, nil)
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
