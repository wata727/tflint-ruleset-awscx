package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsDBInstanceStorageThroughputNonGP3Rule checks that storage_throughput is only used with gp3 storage.
type AwsDBInstanceStorageThroughputNonGP3Rule struct {
	tflint.DefaultRule

	resourceType               string
	storageTypeAttribute       string
	storageThroughputAttribute string
}

// NewAwsDBInstanceStorageThroughputNonGP3Rule returns a new rule.
func NewAwsDBInstanceStorageThroughputNonGP3Rule() *AwsDBInstanceStorageThroughputNonGP3Rule {
	return &AwsDBInstanceStorageThroughputNonGP3Rule{
		resourceType:               "aws_db_instance",
		storageTypeAttribute:       "storage_type",
		storageThroughputAttribute: "storage_throughput",
	}
}

// Name returns the rule name.
func (r *AwsDBInstanceStorageThroughputNonGP3Rule) Name() string {
	return "awscx_db_instance_storage_throughput_non_gp3"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstanceStorageThroughputNonGP3Rule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstanceStorageThroughputNonGP3Rule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstanceStorageThroughputNonGP3Rule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance"
}

// Check reports when storage_throughput is configured without gp3 storage.
func (r *AwsDBInstanceStorageThroughputNonGP3Rule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.storageTypeAttribute},
			{Name: r.storageThroughputAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		storageThroughput, hasStorageThroughput := resource.Body.Attributes[r.storageThroughputAttribute]
		if !hasStorageThroughput {
			continue
		}

		storageType, hasStorageType := resource.Body.Attributes[r.storageTypeAttribute]
		if !hasStorageType {
			runner.EmitIssue(
				r,
				"`storage_throughput` can only be set when `storage_type = \"gp3\"`.",
				storageThroughput.Expr.Range(),
			)
			continue
		}

		err := runner.EvaluateExpr(storageType.Expr, func(value string) error {
			if strings.EqualFold(value, "gp3") {
				return nil
			}

			runner.EmitIssue(
				r,
				"`storage_throughput` can only be set when `storage_type = \"gp3\"`.",
				storageThroughput.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
