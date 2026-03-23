package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLBTargetGroupMatcherNonHTTPHealthCheckRule checks matcher usage on non-HTTP health checks.
type AwsLBTargetGroupMatcherNonHTTPHealthCheckRule struct {
	tflint.DefaultRule

	resourceType        string
	targetTypeAttribute string
	healthCheckBlock    string
	protocolAttribute   string
	matcherAttribute    string
}

// NewAwsLBTargetGroupMatcherNonHTTPHealthCheckRule returns a new rule.
func NewAwsLBTargetGroupMatcherNonHTTPHealthCheckRule() *AwsLBTargetGroupMatcherNonHTTPHealthCheckRule {
	return &AwsLBTargetGroupMatcherNonHTTPHealthCheckRule{
		resourceType:        "aws_lb_target_group",
		targetTypeAttribute: "target_type",
		healthCheckBlock:    "health_check",
		protocolAttribute:   "protocol",
		matcherAttribute:    "matcher",
	}
}

// Name returns the rule name.
func (r *AwsLBTargetGroupMatcherNonHTTPHealthCheckRule) Name() string {
	return "awscx_lb_target_group_matcher_non_http_health_check"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLBTargetGroupMatcherNonHTTPHealthCheckRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLBTargetGroupMatcherNonHTTPHealthCheckRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLBTargetGroupMatcherNonHTTPHealthCheckRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group"
}

// Check reports matcher usage on health checks that are not HTTP or HTTPS.
func (r *AwsLBTargetGroupMatcherNonHTTPHealthCheckRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.targetTypeAttribute},
		},
		Blocks: []hclext.BlockSchema{
			{
				Type: r.healthCheckBlock,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.protocolAttribute},
						{Name: r.matcherAttribute},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		if shouldSkip, err := r.isLambdaTargetGroup(runner, resource); err != nil {
			return err
		} else if shouldSkip {
			continue
		}

		for _, block := range resource.Body.Blocks {
			matcher, hasMatcher := block.Body.Attributes[r.matcherAttribute]
			if !hasMatcher {
				continue
			}

			protocol, hasProtocol := block.Body.Attributes[r.protocolAttribute]
			if !hasProtocol {
				continue
			}

			err := runner.EvaluateExpr(protocol.Expr, func(value string) error {
				switch strings.ToUpper(strings.TrimSpace(value)) {
				case "HTTP", "HTTPS":
					return nil
				}

				runner.EmitIssue(
					r,
					"`health_check.matcher` can only be set when `health_check.protocol` is `HTTP` or `HTTPS`.",
					matcher.Expr.Range(),
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

func (r *AwsLBTargetGroupMatcherNonHTTPHealthCheckRule) isLambdaTargetGroup(runner tflint.Runner, resource *hclext.Block) (bool, error) {
	targetType, exists := resource.Body.Attributes[r.targetTypeAttribute]
	if !exists {
		return false, nil
	}

	isLambda := false
	evaluated := false

	err := runner.EvaluateExpr(targetType.Expr, func(value string) error {
		evaluated = true
		isLambda = strings.EqualFold(strings.TrimSpace(value), "lambda")
		return nil
	}, nil)
	if err != nil {
		return false, err
	}

	if !evaluated {
		return true, nil
	}

	return isLambda, nil
}
