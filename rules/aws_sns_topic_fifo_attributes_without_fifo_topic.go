package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule checks FIFO-only SNS topic arguments are used only with FIFO topics.
type AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule struct {
	tflint.DefaultRule

	resourceType                 string
	fifoTopicAttribute           string
	archivePolicyAttribute       string
	contentBasedDedupAttribute   string
	fifoThroughputScopeAttribute string
}

// NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule returns a new rule.
func NewAwsSNSTopicFIFOAttributesWithoutFIFOTopicRule() *AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule {
	return &AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule{
		resourceType:                 "aws_sns_topic",
		fifoTopicAttribute:           "fifo_topic",
		archivePolicyAttribute:       "archive_policy",
		contentBasedDedupAttribute:   "content_based_deduplication",
		fifoThroughputScopeAttribute: "fifo_throughput_scope",
	}
}

// Name returns the rule name.
func (r *AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule) Name() string {
	return "awscx_sns_topic_fifo_attributes_without_fifo_topic"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic"
}

// Check reports FIFO-only SNS topic arguments used without fifo_topic enabled.
func (r *AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.fifoTopicAttribute},
			{Name: r.archivePolicyAttribute},
			{Name: r.contentBasedDedupAttribute},
			{Name: r.fifoThroughputScopeAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		fifoTopic, hasFIFOTopic := resource.Body.Attributes[r.fifoTopicAttribute]
		archivePolicy, hasArchivePolicy := resource.Body.Attributes[r.archivePolicyAttribute]
		contentBasedDeduplication, hasContentBasedDeduplication := resource.Body.Attributes[r.contentBasedDedupAttribute]
		fifoThroughputScope, hasFIFOThroughputScope := resource.Body.Attributes[r.fifoThroughputScopeAttribute]

		if !hasArchivePolicy && !hasContentBasedDeduplication && !hasFIFOThroughputScope {
			continue
		}

		if !hasFIFOTopic {
			r.emitIssues(runner, hasArchivePolicy, archivePolicy, hasContentBasedDeduplication, contentBasedDeduplication, hasFIFOThroughputScope, fifoThroughputScope)
			continue
		}

		err := runner.EvaluateExpr(fifoTopic.Expr, func(enabled bool) error {
			if enabled {
				return nil
			}

			r.emitIssues(runner, hasArchivePolicy, archivePolicy, hasContentBasedDeduplication, contentBasedDeduplication, hasFIFOThroughputScope, fifoThroughputScope)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *AwsSNSTopicFIFOAttributesWithoutFIFOTopicRule) emitIssues(
	runner tflint.Runner,
	hasArchivePolicy bool,
	archivePolicy *hclext.Attribute,
	hasContentBasedDeduplication bool,
	contentBasedDeduplication *hclext.Attribute,
	hasFIFOThroughputScope bool,
	fifoThroughputScope *hclext.Attribute,
) {
	if hasArchivePolicy {
		runner.EmitIssue(
			r,
			"`archive_policy` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
			archivePolicy.Expr.Range(),
		)
	}

	if hasContentBasedDeduplication {
		runner.EmitIssue(
			r,
			"`content_based_deduplication` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
			contentBasedDeduplication.Expr.Range(),
		)
	}

	if hasFIFOThroughputScope {
		runner.EmitIssue(
			r,
			"`fifo_throughput_scope` cannot be set unless `fifo_topic = true` on `aws_sns_topic`.",
			fifoThroughputScope.Expr.Range(),
		)
	}
}
