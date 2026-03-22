package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedLifecycleRule warns when deprecated inline lifecycle_rule is used.
type AwsS3BucketDeprecatedLifecycleRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedLifecycleRule returns a new rule.
func NewAwsS3BucketDeprecatedLifecycleRule() *AwsS3BucketDeprecatedLifecycleRule {
	return &AwsS3BucketDeprecatedLifecycleRule{
		resourceType: "aws_s3_bucket",
		blockType:    "lifecycle_rule",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedLifecycleRule) Name() string {
	return "awscx_s3_bucket_deprecated_lifecycle_rule"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedLifecycleRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedLifecycleRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedLifecycleRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline lifecycle_rule is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedLifecycleRule) Check(runner tflint.Runner) error {
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
				"`lifecycle_rule` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_lifecycle_configuration` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
