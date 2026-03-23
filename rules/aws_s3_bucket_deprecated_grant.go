package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedGrantRule warns when deprecated inline grant is used.
type AwsS3BucketDeprecatedGrantRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedGrantRule returns a new rule.
func NewAwsS3BucketDeprecatedGrantRule() *AwsS3BucketDeprecatedGrantRule {
	return &AwsS3BucketDeprecatedGrantRule{
		resourceType: "aws_s3_bucket",
		blockType:    "grant",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedGrantRule) Name() string {
	return "awscx_s3_bucket_deprecated_grant"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedGrantRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedGrantRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedGrantRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline grant blocks are configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedGrantRule) Check(runner tflint.Runner) error {
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
				"`grant` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_acl` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
