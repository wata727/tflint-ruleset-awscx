package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketDeprecatedACLRule warns when deprecated inline acl is used.
type AwsS3BucketDeprecatedACLRule struct {
	tflint.DefaultRule

	resourceType  string
	attributeName string
}

// NewAwsS3BucketDeprecatedACLRule returns a new rule.
func NewAwsS3BucketDeprecatedACLRule() *AwsS3BucketDeprecatedACLRule {
	return &AwsS3BucketDeprecatedACLRule{
		resourceType:  "aws_s3_bucket",
		attributeName: "acl",
	}
}

// Name returns the rule name.
func (r *AwsS3BucketDeprecatedACLRule) Name() string {
	return "awscx_s3_bucket_deprecated_acl"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketDeprecatedACLRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketDeprecatedACLRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsS3BucketDeprecatedACLRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket"
}

// Check warns when deprecated inline acl is configured on aws_s3_bucket.
func (r *AwsS3BucketDeprecatedACLRule) Check(runner tflint.Runner) error {
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
			"`acl` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_acl` and `aws_s3_bucket_ownership_controls` when ACLs are required.",
			attribute.Expr.Range(),
		)
	}

	return nil
}
