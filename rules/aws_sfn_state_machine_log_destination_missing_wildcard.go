package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsSFNStateMachineLogDestinationMissingWildcardRule checks Step Functions log destinations use the required :* suffix.
type AwsSFNStateMachineLogDestinationMissingWildcardRule struct {
	tflint.DefaultRule

	resourceType string
	blockType    string
	attribute    string
}

// NewAwsSFNStateMachineLogDestinationMissingWildcardRule returns a new rule.
func NewAwsSFNStateMachineLogDestinationMissingWildcardRule() *AwsSFNStateMachineLogDestinationMissingWildcardRule {
	return &AwsSFNStateMachineLogDestinationMissingWildcardRule{
		resourceType: "aws_sfn_state_machine",
		blockType:    "logging_configuration",
		attribute:    "log_destination",
	}
}

// Name returns the rule name.
func (r *AwsSFNStateMachineLogDestinationMissingWildcardRule) Name() string {
	return "awscx_sfn_state_machine_log_destination_missing_wildcard"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsSFNStateMachineLogDestinationMissingWildcardRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsSFNStateMachineLogDestinationMissingWildcardRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsSFNStateMachineLogDestinationMissingWildcardRule) Link() string {
	return "https://docs.aws.amazon.com/step-functions/latest/apireference/API_CloudWatchLogsLogGroup.html"
}

// Check reports Step Functions log destinations that omit the required CloudWatch Logs wildcard suffix.
func (r *AwsSFNStateMachineLogDestinationMissingWildcardRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: r.attribute}},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			attribute, exists := block.Body.Attributes[r.attribute]
			if !exists {
				continue
			}

			err := runner.EvaluateExpr(attribute.Expr, func(value cty.Value) error {
				if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
					return nil
				}

				if strings.HasSuffix(strings.TrimSpace(value.AsString()), ":*") {
					return nil
				}

				runner.EmitIssue(
					r,
					"`logging_configuration.log_destination` must end with `:*` for Step Functions CloudWatch Logs destinations.",
					attribute.Expr.Range(),
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
