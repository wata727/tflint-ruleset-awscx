package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsLBTargetGroupProtocolVersionNonHTTPRule checks protocol_version usage on non-HTTP target groups.
type AwsLBTargetGroupProtocolVersionNonHTTPRule struct {
	tflint.DefaultRule

	resourceType             string
	protocolAttribute        string
	protocolVersionAttribute string
	targetTypeAttribute      string
}

// NewAwsLBTargetGroupProtocolVersionNonHTTPRule returns a new rule.
func NewAwsLBTargetGroupProtocolVersionNonHTTPRule() *AwsLBTargetGroupProtocolVersionNonHTTPRule {
	return &AwsLBTargetGroupProtocolVersionNonHTTPRule{
		resourceType:             "aws_lb_target_group",
		protocolAttribute:        "protocol",
		protocolVersionAttribute: "protocol_version",
		targetTypeAttribute:      "target_type",
	}
}

// Name returns the rule name.
func (r *AwsLBTargetGroupProtocolVersionNonHTTPRule) Name() string {
	return "awscx_lb_target_group_protocol_version_non_http"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLBTargetGroupProtocolVersionNonHTTPRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLBTargetGroupProtocolVersionNonHTTPRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLBTargetGroupProtocolVersionNonHTTPRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group"
}

// Check reports target groups that set protocol_version outside HTTP/HTTPS target groups.
func (r *AwsLBTargetGroupProtocolVersionNonHTTPRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.protocolAttribute},
			{Name: r.protocolVersionAttribute},
			{Name: r.targetTypeAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		protocolVersion, hasProtocolVersion := resource.Body.Attributes[r.protocolVersionAttribute]
		if !hasProtocolVersion {
			continue
		}

		if isLambda, known, err := r.isLambdaTargetGroup(runner, resource); err != nil {
			return err
		} else if known && isLambda {
			runner.EmitIssue(
				r,
				"`protocol_version` can only be set when `protocol` is `HTTP` or `HTTPS`.",
				protocolVersion.Expr.Range(),
			)
			continue
		}

		protocol, hasProtocol := resource.Body.Attributes[r.protocolAttribute]
		if !hasProtocol {
			continue
		}

		err := runner.EvaluateExpr(protocol.Expr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			switch strings.ToUpper(strings.TrimSpace(value.AsString())) {
			case "HTTP", "HTTPS":
				return nil
			default:
				runner.EmitIssue(
					r,
					"`protocol_version` can only be set when `protocol` is `HTTP` or `HTTPS`.",
					protocolVersion.Expr.Range(),
				)
				return nil
			}
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *AwsLBTargetGroupProtocolVersionNonHTTPRule) isLambdaTargetGroup(runner tflint.Runner, resource *hclext.Block) (bool, bool, error) {
	targetType, exists := resource.Body.Attributes[r.targetTypeAttribute]
	if !exists {
		return false, false, nil
	}

	isLambda := false
	evaluated := false

	err := runner.EvaluateExpr(targetType.Expr, func(value cty.Value) error {
		if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
			return nil
		}

		evaluated = true
		isLambda = strings.EqualFold(strings.TrimSpace(value.AsString()), "lambda")
		return nil
	}, nil)
	if err != nil {
		return false, false, err
	}

	return isLambda, evaluated, nil
}
