package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedRequestPayerRule warns when deprecated inline request_payer is used.
type AwsS3BucketDeprecatedRequestPayerRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsS3BucketDeprecatedRequestPayerRule returns a new rule.
func NewAwsS3BucketDeprecatedRequestPayerRule() *AwsS3BucketDeprecatedRequestPayerRule {
	return &AwsS3BucketDeprecatedRequestPayerRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "request_payer",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedRequestPayerRule) Name() string {
	return "awscx_s3_bucket_deprecated_request_payer"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedRequestPayerRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedRequestPayerRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedRequestPayerRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline request_payer is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedRequestPayerRule) Check(runner tflint.Runner) error {
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
			"`request_payer` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_request_payment_configuration` instead.",
			attribute.Expr.Range(),
		)
	}

	return nil
}
