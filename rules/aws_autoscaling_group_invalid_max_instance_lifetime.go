package rules

import (
	"math/big"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule checks max_instance_lifetime values on aws_autoscaling_group.
type AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule struct {
	tflint.DefaultRule

	resourceType             string
	maxInstanceLifetimeAttr  string
	minimumLifetimeInSeconds int64
	maximumLifetimeInSeconds int64
}

// NewAwsAutoscalingGroupInvalidMaxInstanceLifetimeRule returns a new rule.
func NewAwsAutoscalingGroupInvalidMaxInstanceLifetimeRule() *AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule {
	return &AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule{
		resourceType:             "aws_autoscaling_group",
		maxInstanceLifetimeAttr:  "max_instance_lifetime",
		minimumLifetimeInSeconds: 86400,
		maximumLifetimeInSeconds: 31536000,
	}
}

// Name returns the rule name.
func (r *AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule) Name() string {
	return "awscx_autoscaling_group_invalid_max_instance_lifetime"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/autoscaling_group"
}

// Check reports explicit max_instance_lifetime values that fall outside the documented range.
func (r *AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.maxInstanceLifetimeAttr},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		maxInstanceLifetime, exists := resource.Body.Attributes[r.maxInstanceLifetimeAttr]
		if !exists {
			continue
		}

		value, ok, err := r.evaluateInteger(runner, maxInstanceLifetime.Expr)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		if value == 0 {
			continue
		}

		if value >= r.minimumLifetimeInSeconds && value <= r.maximumLifetimeInSeconds {
			continue
		}

		runner.EmitIssue(
			r,
			"`max_instance_lifetime` must be 0 or between 86400 and 31536000 seconds on `aws_autoscaling_group`.",
			maxInstanceLifetime.Expr.Range(),
		)
	}

	return nil
}

func (r *AwsAutoscalingGroupInvalidMaxInstanceLifetimeRule) evaluateInteger(runner tflint.Runner, expr hcl.Expression) (int64, bool, error) {
	var value int64
	resolved := false

	err := runner.EvaluateExpr(expr, func(result cty.Value) error {
		if !result.IsKnown() || result.IsNull() || !result.Type().Equals(cty.Number) {
			return nil
		}

		intValue, accuracy := result.AsBigFloat().Int64()
		if accuracy != big.Exact {
			return nil
		}

		value = intValue
		resolved = true
		return nil
	}, nil)
	if err != nil {
		return 0, false, err
	}

	return value, resolved, nil
}
