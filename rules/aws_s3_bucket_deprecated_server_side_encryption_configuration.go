package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule warns on deprecated inline bucket encryption configuration.
type AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedServerSideEncryptionConfigurationRule returns a new rule.
func NewAwsS3BucketDeprecatedServerSideEncryptionConfigurationRule() *AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule {
	return &AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule{
		resourceType: "aws_s3_bucket",
		blockType:    "server_side_encryption_configuration",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule) Name() string {
	return "awscx_s3_bucket_deprecated_server_side_encryption_configuration"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline bucket encryption configuration is used.
func (r *AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule) Check(runner tflint.Runner) error {
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
				"`server_side_encryption_configuration` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_server_side_encryption_configuration` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
