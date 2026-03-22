package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedReplicationConfigurationRule warns when deprecated inline replication_configuration is used.
type AwsS3BucketDeprecatedReplicationConfigurationRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
}

// NewAwsS3BucketDeprecatedReplicationConfigurationRule returns a new rule.
func NewAwsS3BucketDeprecatedReplicationConfigurationRule() *AwsS3BucketDeprecatedReplicationConfigurationRule {
	return &AwsS3BucketDeprecatedReplicationConfigurationRule{
		resourceType: "aws_s3_bucket",
		blockType:    "replication_configuration",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedReplicationConfigurationRule) Name() string {
	return "awscx_s3_bucket_deprecated_replication_configuration"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedReplicationConfigurationRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedReplicationConfigurationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedReplicationConfigurationRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline replication_configuration is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedReplicationConfigurationRule) Check(runner tflint.Runner) error {
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
				"`replication_configuration` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_replication_configuration` instead.",
				block.DefRange,
			)
		}
	}

	return nil
}
