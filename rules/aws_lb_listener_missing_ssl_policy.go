package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLBListenerMissingSSLPolicyRule checks HTTPS/TLS listeners declare an SSL policy.
type AwsLBListenerMissingSSLPolicyRule struct {
	tflint.DefaultRule

	resourceType       string
	protocolAttribute  string
	sslPolicyAttribute string
}

// NewAwsLBListenerMissingSSLPolicyRule returns a new rule.
func NewAwsLBListenerMissingSSLPolicyRule() *AwsLBListenerMissingSSLPolicyRule {
	return &AwsLBListenerMissingSSLPolicyRule{
		resourceType:       "aws_lb_listener",
		protocolAttribute:  "protocol",
		sslPolicyAttribute: "ssl_policy",
	}
}

// Name returns the rule name.
func (r *AwsLBListenerMissingSSLPolicyRule) Name() string {
	return "awscx_lb_listener_missing_ssl_policy"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLBListenerMissingSSLPolicyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLBListenerMissingSSLPolicyRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLBListenerMissingSSLPolicyRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener"
}

// Check reports HTTPS/TLS listeners that omit the required SSL policy.
func (r *AwsLBListenerMissingSSLPolicyRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.protocolAttribute},
			{Name: r.sslPolicyAttribute},
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

		_, hasSSLPolicy := resource.Body.Attributes[r.sslPolicyAttribute]

		err := runner.EvaluateExpr(protocol.Expr, func(value string) error {
			switch strings.ToUpper(strings.TrimSpace(value)) {
			case "HTTPS", "TLS":
			default:
				return nil
			}

			if hasSSLPolicy {
				return nil
			}

			runner.EmitIssue(
				r,
				"`ssl_policy` must be set when `protocol` is `HTTPS` or `TLS`.",
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
