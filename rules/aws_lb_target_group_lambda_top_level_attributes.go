package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsLBTargetGroupLambdaTopLevelAttributesRule checks top-level attributes that do not apply to Lambda target groups.
type AwsLBTargetGroupLambdaTopLevelAttributesRule struct {
	tflint.DefaultRule

	resourceType        string
	targetTypeAttribute string
	invalidAttributes   []string
}

// NewAwsLBTargetGroupLambdaTopLevelAttributesRule returns a new rule.
func NewAwsLBTargetGroupLambdaTopLevelAttributesRule() *AwsLBTargetGroupLambdaTopLevelAttributesRule {
	return &AwsLBTargetGroupLambdaTopLevelAttributesRule{
		resourceType:        "aws_lb_target_group",
		targetTypeAttribute: "target_type",
		invalidAttributes:   []string{"port", "protocol", "vpc_id"},
	}
}

// Name returns the rule name.
func (r *AwsLBTargetGroupLambdaTopLevelAttributesRule) Name() string {
	return "awscx_lb_target_group_lambda_top_level_attributes"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLBTargetGroupLambdaTopLevelAttributesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLBTargetGroupLambdaTopLevelAttributesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLBTargetGroupLambdaTopLevelAttributesRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group"
}

// Check reports port, protocol, and vpc_id on explicit Lambda target groups.
func (r *AwsLBTargetGroupLambdaTopLevelAttributesRule) Check(runner tflint.Runner) error {
	schema := &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.targetTypeAttribute},
		},
	}
	for _, name := range r.invalidAttributes {
		schema.Attributes = append(schema.Attributes, hclext.AttributeSchema{Name: name})
	}

	resources, err := runner.GetResourceContent(r.resourceType, schema, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		targetType, exists := resource.Body.Attributes[r.targetTypeAttribute]
		if !exists {
			continue
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
			return err
		}

		if !evaluated || !isLambda {
			continue
		}

		for _, name := range r.invalidAttributes {
			attribute, exists := resource.Body.Attributes[name]
			if !exists {
				continue
			}

			runner.EmitIssue(
				r,
				"`"+name+"` does not apply when `target_type = \"lambda\"`.",
				attribute.Expr.Range(),
			)
		}
	}

	return nil
}
