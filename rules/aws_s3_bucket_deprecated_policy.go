package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedPolicyRule warns when deprecated inline policy is used.
type AwsS3BucketDeprecatedPolicyRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsS3BucketDeprecatedPolicyRule returns a new rule.
func NewAwsS3BucketDeprecatedPolicyRule() *AwsS3BucketDeprecatedPolicyRule {
	return &AwsS3BucketDeprecatedPolicyRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "policy",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedPolicyRule) Name() string {
	return "awscx_s3_bucket_deprecated_policy"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedPolicyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedPolicyRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedPolicyRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline policy is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedPolicyRule) Check(runner tflint.Runner) error {
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
			"`policy` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_policy` instead.",
			attribute.Expr.Range(),
		)
	}

	return nil
}
