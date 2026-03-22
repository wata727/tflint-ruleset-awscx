package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSQSQueueFIFONameSuffixRule checks FIFO queue names use the required suffix.
type AwsSQSQueueFIFONameSuffixRule struct {
	tflint.DefaultRule

	resourceType       string
	nameAttribute      string
	fifoQueueAttribute string
	requiredNameSuffix string
}

// NewAwsSQSQueueFIFONameSuffixRule returns a new rule.
func NewAwsSQSQueueFIFONameSuffixRule() *AwsSQSQueueFIFONameSuffixRule {
	return &AwsSQSQueueFIFONameSuffixRule{
		resourceType:       "aws_sqs_queue",
		nameAttribute:      "name",
		fifoQueueAttribute: "fifo_queue",
		requiredNameSuffix: ".fifo",
	}
}

// Name returns the rule name.
func (r *AwsSQSQueueFIFONameSuffixRule) Name() string {
	return "awscx_sqs_queue_fifo_name_suffix"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsSQSQueueFIFONameSuffixRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsSQSQueueFIFONameSuffixRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsSQSQueueFIFONameSuffixRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue"
}

// Check reports FIFO queues whose explicit name omits the required .fifo suffix.
func (r *AwsSQSQueueFIFONameSuffixRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.nameAttribute},
			{Name: r.fifoQueueAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		name, hasName := resource.Body.Attributes[r.nameAttribute]
		fifoQueue, hasFIFOQueue := resource.Body.Attributes[r.fifoQueueAttribute]
		if !hasName || !hasFIFOQueue {
			continue
		}

		err := runner.EvaluateExpr(fifoQueue.Expr, func(enabled bool) error {
			if !enabled {
				return nil
			}

			return runner.EvaluateExpr(name.Expr, func(queueName string) error {
				if strings.HasSuffix(queueName, r.requiredNameSuffix) {
					return nil
				}

				runner.EmitIssue(
					r,
					"`name` must end with `.fifo` when `fifo_queue = true`.",
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
