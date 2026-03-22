package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsDBInstanceMissingIOPSRule checks RDS storage settings that require IOPS.
type AwsDBInstanceMissingIOPSRule struct {
	tflint.DefaultRule

	resourceType string
	storageType  string
	iops         string
}

// NewAwsDBInstanceMissingIOPSRule returns a new rule.
func NewAwsDBInstanceMissingIOPSRule() *AwsDBInstanceMissingIOPSRule {
	return &AwsDBInstanceMissingIOPSRule{
		resourceType: "aws_db_instance",
		storageType:  "storage_type",
		iops:         "iops",
	}
}

// Name returns the rule name.
func (r *AwsDBInstanceMissingIOPSRule) Name() string {
	return "awscx_db_instance_missing_iops"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstanceMissingIOPSRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstanceMissingIOPSRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstanceMissingIOPSRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance"
}

// Check reports when storage types that require IOPS omit it.
func (r *AwsDBInstanceMissingIOPSRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.storageType},
			{Name: r.iops},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		storageType, exists := resource.Body.Attributes[r.storageType]
		if !exists {
			continue
		}

		_, hasIOPS := resource.Body.Attributes[r.iops]

		err := runner.EvaluateExpr(storageType.Expr, func(value string) error {
			switch strings.ToLower(value) {
			case "io1", "io2", "gp3":
			default:
				return nil
			}

			if hasIOPS {
				return nil
			}

			runner.EmitIssue(
				r,
				"`iops` must be set when `storage_type` is `io1`, `io2`, or `gp3`.",
				storageType.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
