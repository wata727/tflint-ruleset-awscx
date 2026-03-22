package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSNSTopicFIFONameSuffixRule checks FIFO topic names use the required suffix.
type AwsSNSTopicFIFONameSuffixRule struct {
	tflint.DefaultRule

	resourceType       string
	nameAttribute      string
	fifoTopicAttribute string
	requiredNameSuffix string
}

// NewAwsSNSTopicFIFONameSuffixRule returns a new rule.
func NewAwsSNSTopicFIFONameSuffixRule() *AwsSNSTopicFIFONameSuffixRule {
	return &AwsSNSTopicFIFONameSuffixRule{
		resourceType:       "aws_sns_topic",
		nameAttribute:      "name",
		fifoTopicAttribute: "fifo_topic",
		requiredNameSuffix: ".fifo",
	}
}

// Name returns the rule name.
func (r *AwsSNSTopicFIFONameSuffixRule) Name() string {
	return "awscx_sns_topic_fifo_name_suffix"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsSNSTopicFIFONameSuffixRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsSNSTopicFIFONameSuffixRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsSNSTopicFIFONameSuffixRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic"
}

// Check reports FIFO topics whose explicit name omits the required .fifo suffix.
func (r *AwsSNSTopicFIFONameSuffixRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.nameAttribute},
			{Name: r.fifoTopicAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		name, hasName := resource.Body.Attributes[r.nameAttribute]
		fifoTopic, hasFIFOTopic := resource.Body.Attributes[r.fifoTopicAttribute]
		if !hasName || !hasFIFOTopic {
			continue
		}

		err := runner.EvaluateExpr(fifoTopic.Expr, func(enabled bool) error {
			if !enabled {
				return nil
			}

			return runner.EvaluateExpr(name.Expr, func(topicName string) error {
				if strings.HasSuffix(topicName, r.requiredNameSuffix) {
					return nil
				}

				runner.EmitIssue(
					r,
					"`name` must end with `.fifo` when `fifo_topic = true`.",
					name.Expr.Range(),
				)
				return nil
			}, nil)
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
