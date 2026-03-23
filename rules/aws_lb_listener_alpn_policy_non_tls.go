package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsLBListenerALPNPolicyNonTLSRule checks that alpn_policy is only used with TLS listeners.
type AwsLBListenerALPNPolicyNonTLSRule struct {
	tflint.DefaultRule

	resourceType         string
	protocolAttribute    string
	alpnPolicyAttribute  string
	requiredProtocolName string
}

// NewAwsLBListenerALPNPolicyNonTLSRule returns a new rule.
func NewAwsLBListenerALPNPolicyNonTLSRule() *AwsLBListenerALPNPolicyNonTLSRule {
	return &AwsLBListenerALPNPolicyNonTLSRule{
		resourceType:         "aws_lb_listener",
		protocolAttribute:    "protocol",
		alpnPolicyAttribute:  "alpn_policy",
		requiredProtocolName: "TLS",
	}
}

// Name returns the rule name.
func (r *AwsLBListenerALPNPolicyNonTLSRule) Name() string {
	return "awscx_lb_listener_alpn_policy_non_tls"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLBListenerALPNPolicyNonTLSRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLBListenerALPNPolicyNonTLSRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLBListenerALPNPolicyNonTLSRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener"
}

// Check reports listeners that set alpn_policy with a non-TLS protocol.
func (r *AwsLBListenerALPNPolicyNonTLSRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.protocolAttribute},
			{Name: r.alpnPolicyAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		protocol, hasProtocol := resource.Body.Attributes[r.protocolAttribute]
		alpnPolicy, hasALPNPolicy := resource.Body.Attributes[r.alpnPolicyAttribute]
		if !hasProtocol || !hasALPNPolicy {
			continue
		}

		err := runner.EvaluateExpr(protocol.Expr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			if strings.EqualFold(strings.TrimSpace(value.AsString()), r.requiredProtocolName) {
				return nil
			}

			runner.EmitIssue(
				r,
				"`alpn_policy` can only be set when `protocol` is `TLS`.",
				alpnPolicy.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
