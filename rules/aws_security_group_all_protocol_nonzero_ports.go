package rules

import (
	"math/big"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsSecurityGroupAllProtocolNonzeroPortsRule checks that all-protocol rules use 0/0 ports.
type AwsSecurityGroupAllProtocolNonzeroPortsRule struct {
	tflint.DefaultRule

	resourceType string
	ingressType  string
	egressType   string
	protocolAttr string
	fromPortAttr string
	toPortAttr   string
}

// NewAwsSecurityGroupAllProtocolNonzeroPortsRule returns a new rule.
func NewAwsSecurityGroupAllProtocolNonzeroPortsRule() *AwsSecurityGroupAllProtocolNonzeroPortsRule {
	return &AwsSecurityGroupAllProtocolNonzeroPortsRule{
		resourceType: "aws_security_group",
		ingressType:  "ingress",
		egressType:   "egress",
		protocolAttr: "protocol",
		fromPortAttr: "from_port",
		toPortAttr:   "to_port",
	}
}

// Name returns the rule name.
func (r *AwsSecurityGroupAllProtocolNonzeroPortsRule) Name() string {
	return "awscx_security_group_all_protocol_nonzero_ports"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsSecurityGroupAllProtocolNonzeroPortsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsSecurityGroupAllProtocolNonzeroPortsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsSecurityGroupAllProtocolNonzeroPortsRule) Link() string {
	return "https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html"
}

// Check reports all-protocol security group rules that do not use 0/0 ports.
func (r *AwsSecurityGroupAllProtocolNonzeroPortsRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.ingressType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.protocolAttr},
						{Name: r.fromPortAttr},
						{Name: r.toPortAttr},
					},
				},
			},
			{
				Type: r.egressType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.protocolAttr},
						{Name: r.fromPortAttr},
						{Name: r.toPortAttr},
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
			protocolAttr, exists := block.Body.Attributes[r.protocolAttr]
			if !exists {
				continue
			}

			fromPortAttr, hasFromPort := block.Body.Attributes[r.fromPortAttr]
			toPortAttr, hasToPort := block.Body.Attributes[r.toPortAttr]
			if !hasFromPort || !hasToPort {
				continue
			}

			err := runner.EvaluateExpr(protocolAttr.Expr, func(protocol string) error {
				if strings.TrimSpace(protocol) != "-1" {
					return nil
				}

				fromPort, ok, err := r.evaluatePort(runner, fromPortAttr.Expr)
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}

				toPort, ok, err := r.evaluatePort(runner, toPortAttr.Expr)
				if err != nil {
					return err
				}
				if !ok {
					return nil
				}

				if fromPort == 0 && toPort == 0 {
					return nil
				}

				runner.EmitIssue(
					r,
					"`from_port` and `to_port` must both be `0` when `protocol = \"-1\"`.",
					protocolAttr.Expr.Range(),
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

func (r *AwsSecurityGroupAllProtocolNonzeroPortsRule) evaluatePort(runner tflint.Runner, expr hcl.Expression) (int, bool, error) {
	port := 0
	resolved := false

	err := runner.EvaluateExpr(expr, func(value cty.Value) error {
		if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.Number) {
			return nil
		}

		portFloat := value.AsBigFloat()
		portInt, accuracy := portFloat.Int64()
		if accuracy != big.Exact {
			return nil
		}

		port = int(portInt)
		resolved = true
		return nil
	}, nil)
	if err != nil {
		return 0, false, err
	}

	return port, resolved, nil
}
