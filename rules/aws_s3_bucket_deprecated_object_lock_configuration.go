package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedObjectLockConfigurationRule warns when deprecated inline object lock configuration is used.
type AwsS3BucketDeprecatedObjectLockConfigurationRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedObjectLockConfigurationRule returns a new rule.
func NewAwsS3BucketDeprecatedObjectLockConfigurationRule() *AwsS3BucketDeprecatedObjectLockConfigurationRule {
	return &AwsS3BucketDeprecatedObjectLockConfigurationRule{
		resourceType: "aws_s3_bucket",
		blockType:    "object_lock_configuration",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedObjectLockConfigurationRule) Name() string {
	return "awscx_s3_bucket_deprecated_object_lock_configuration"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedObjectLockConfigurationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedObjectLockConfigurationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedObjectLockConfigurationRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline object_lock_configuration is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedObjectLockConfigurationRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: r.blockType},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			if block.Type != r.blockType {
				continue
			}

			runner.EmitIssue(
				r,
				"`object_lock_configuration` on `aws_s3_bucket` is deprecated; use `object_lock_enabled` with `aws_s3_bucket_object_lock_configuration` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
