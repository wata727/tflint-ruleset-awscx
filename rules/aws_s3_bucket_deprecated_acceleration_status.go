package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedAccelerationStatusRule warns when deprecated inline acceleration_status is used.
type AwsS3BucketDeprecatedAccelerationStatusRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsS3BucketDeprecatedAccelerationStatusRule returns a new rule.
func NewAwsS3BucketDeprecatedAccelerationStatusRule() *AwsS3BucketDeprecatedAccelerationStatusRule {
	return &AwsS3BucketDeprecatedAccelerationStatusRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "acceleration_status",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedAccelerationStatusRule) Name() string {
	return "awscx_s3_bucket_deprecated_acceleration_status"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedAccelerationStatusRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedAccelerationStatusRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedAccelerationStatusRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline acceleration_status is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedAccelerationStatusRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attributeName}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attributeName]
		if !exists {
			continue
		}

		runner.EmitIssue(
			r,
			"`acceleration_status` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_accelerate_configuration` instead.",
			attribute.Expr.Range(),
		)
	}

	return nil
}
