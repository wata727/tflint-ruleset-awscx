package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedVersioningRule warns when deprecated inline versioning is used.
type AwsS3BucketDeprecatedVersioningRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedVersioningRule returns a new rule.
func NewAwsS3BucketDeprecatedVersioningRule() *AwsS3BucketDeprecatedVersioningRule {
	return &AwsS3BucketDeprecatedVersioningRule{
		resourceType: "aws_s3_bucket",
		blockType:    "versioning",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedVersioningRule) Name() string {
	return "awscx_s3_bucket_deprecated_versioning"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedVersioningRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedVersioningRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedVersioningRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline versioning is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedVersioningRule) Check(runner tflint.Runner) error {
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
				"`versioning` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_versioning` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
