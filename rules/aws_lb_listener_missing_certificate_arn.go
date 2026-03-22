package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLBListenerMissingCertificateARNRule checks HTTPS/TLS listeners declare a default certificate.
type AwsLBListenerMissingCertificateARNRule struct {
	tflint.DefaultRule

	resourceType         string
	protocolAttribute    string
	certificateAttribute string
}

// NewAwsLBListenerMissingCertificateARNRule returns a new rule.
func NewAwsLBListenerMissingCertificateARNRule() *AwsLBListenerMissingCertificateARNRule {
	return &AwsLBListenerMissingCertificateARNRule{
		resourceType:         "aws_lb_listener",
		protocolAttribute:    "protocol",
		certificateAttribute: "certificate_arn",
	}
}

// Name returns the rule name.
func (r *AwsLBListenerMissingCertificateARNRule) Name() string {
	return "awscx_lb_listener_missing_certificate_arn"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLBListenerMissingCertificateARNRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLBListenerMissingCertificateARNRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLBListenerMissingCertificateARNRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener"
}

// Check reports HTTPS/TLS listeners that omit the default certificate ARN.
func (r *AwsLBListenerMissingCertificateARNRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.protocolAttribute},
			{Name: r.certificateAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		protocol, hasProtocol := resource.Body.Attributes[r.protocolAttribute]
		if !hasProtocol {
			continue
		}

		_, hasCertificateARN := resource.Body.Attributes[r.certificateAttribute]

		err := runner.EvaluateExpr(protocol.Expr, func(value string) error {
			switch strings.ToUpper(strings.TrimSpace(value)) {
			case "HTTPS", "TLS":
			default:
				return nil
			}

			if hasCertificateARN {
				return nil
			}

			runner.EmitIssue(
				r,
				"`certificate_arn` must be set when `protocol` is `HTTPS` or `TLS`.",
				protocol.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
