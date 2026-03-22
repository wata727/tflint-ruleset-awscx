package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedWebsiteRule warns when deprecated inline website is used.
type AwsS3BucketDeprecatedWebsiteRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedWebsiteRule returns a new rule.
func NewAwsS3BucketDeprecatedWebsiteRule() *AwsS3BucketDeprecatedWebsiteRule {
	return &AwsS3BucketDeprecatedWebsiteRule{
		resourceType: "aws_s3_bucket",
		blockType:    "website",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedWebsiteRule) Name() string {
	return "awscx_s3_bucket_deprecated_website"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedWebsiteRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedWebsiteRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedWebsiteRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline website is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedWebsiteRule) Check(runner tflint.Runner) error {
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
				"`website` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_website_configuration` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
