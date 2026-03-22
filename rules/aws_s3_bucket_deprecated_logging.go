package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedLoggingRule warns when deprecated inline logging is used.
type AwsS3BucketDeprecatedLoggingRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedLoggingRule returns a new rule.
func NewAwsS3BucketDeprecatedLoggingRule() *AwsS3BucketDeprecatedLoggingRule {
	return &AwsS3BucketDeprecatedLoggingRule{
		resourceType: "aws_s3_bucket",
		blockType:    "logging",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedLoggingRule) Name() string {
	return "awscx_s3_bucket_deprecated_logging"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedLoggingRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedLoggingRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedLoggingRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline logging is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedLoggingRule) Check(runner tflint.Runner) error {
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
				"`logging` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_logging` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
