package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule checks invalid viewer_certificate combinations.
type AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule struct {
	tflint.DefaultRule

	resourceType                    string
	viewerCertificateBlock          string
	defaultCertificateAttribute     string
	minimumProtocolVersionAttribute string
}

// NewAwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule returns a new rule.
func NewAwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule() *AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule {
	return &AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule{
		resourceType:                    "aws_cloudfront_distribution",
		viewerCertificateBlock:          "viewer_certificate",
		defaultCertificateAttribute:     "cloudfront_default_certificate",
		minimumProtocolVersionAttribute: "minimum_protocol_version",
	}
}

// Name returns the rule name.
func (r *AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule) Name() string {
	return "awscx_cloudfront_distribution_minimum_protocol_version_default_certificate"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudfront_distribution"
}

// Check reports minimum_protocol_version usage with the default CloudFront certificate.
func (r *AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.viewerCertificateBlock,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.defaultCertificateAttribute},
						{Name: r.minimumProtocolVersionAttribute},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			minimumProtocolVersion, hasMinimumProtocolVersion := block.Body.Attributes[r.minimumProtocolVersionAttribute]
			if !hasMinimumProtocolVersion {
				continue
			}

			defaultCertificate, hasDefaultCertificate := block.Body.Attributes[r.defaultCertificateAttribute]
			if !hasDefaultCertificate {
				continue
			}

			err := runner.EvaluateExpr(defaultCertificate.Expr, func(value cty.Value) error {
				if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.Bool) {
					return nil
				}

				if !value.True() {
					return nil
				}

				runner.EmitIssue(
					r,
					"`viewer_certificate.minimum_protocol_version` cannot be set when `viewer_certificate.cloudfront_default_certificate` is `true`.",
					minimumProtocolVersion.Expr.Range(),
				)
				return nil
			}, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
