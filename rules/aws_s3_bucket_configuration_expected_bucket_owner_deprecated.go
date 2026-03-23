package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule warns when deprecated expected_bucket_owner is used on S3 bucket configuration sub-resources.
type AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule struct {
	tflint.DefaultRule

	resourceTypes []string
	attributeName string
}

// NewAwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule returns a new rule.
func NewAwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule() *AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule {
	return &AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule{
		resourceTypes: []string{
			"aws_s3_bucket_abac",
			"aws_s3_bucket_accelerate_configuration",
			"aws_s3_bucket_acl",
			"aws_s3_bucket_cors_configuration",
			"aws_s3_bucket_lifecycle_configuration",
			"aws_s3_bucket_logging",
			"aws_s3_bucket_metadata_configuration",
			"aws_s3_bucket_object_lock_configuration",
			"aws_s3_bucket_request_payment_configuration",
			"aws_s3_bucket_server_side_encryption_configuration",
			"aws_s3_bucket_versioning",
			"aws_s3_bucket_website_configuration",
		},
		attributeName: "expected_bucket_owner",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule) Name() string {
	return "awscx_s3_bucket_configuration_expected_bucket_owner_deprecated"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule) Link() string {
	return "https://github.com/hashicorp/terraform-provider-aws/pull/46262"
}

// Check warns when deprecated expected_bucket_owner is configured on S3 bucket configuration sub-resources.
func (r *AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule) Check(runner tflint.Runner) error {
	for _, resourceType := range r.resourceTypes {
		resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{
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
				fmt.Sprintf("`expected_bucket_owner` on `%s` is deprecated; remove it from configuration.", resourceType),
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
