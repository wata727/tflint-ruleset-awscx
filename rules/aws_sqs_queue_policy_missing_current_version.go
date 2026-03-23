package rules

import (
	"encoding/json"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsSQSQueuePolicyMissingCurrentVersionRule checks SQS queue policies explicitly set the current IAM policy version.
type AwsSQSQueuePolicyMissingCurrentVersionRule struct {
	tflint.DefaultRule

	resourceType     string
	policyAttribute  string
	requiredVersion  string
	documentationURL string
}

// NewAwsSQSQueuePolicyMissingCurrentVersionRule returns a new rule.
func NewAwsSQSQueuePolicyMissingCurrentVersionRule() *AwsSQSQueuePolicyMissingCurrentVersionRule {
	return &AwsSQSQueuePolicyMissingCurrentVersionRule{
		resourceType:     "aws_sqs_queue_policy",
		policyAttribute:  "policy",
		requiredVersion:  "2012-10-17",
		documentationURL: "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue",
	}
}

// Name returns the rule name.
func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) Name() string {
	return "awscx_sqs_queue_policy_missing_current_version"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) Link() string {
	return r.documentationURL
}

// Check reports SQS queue policies that omit or override the required top-level Version field.
func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.policyAttribute}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		policy, exists := resource.Body.Attributes[r.policyAttribute]
		if !exists {
			continue
		}

		version, ok, err := r.extractVersion(runner, policy.Expr)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		if version == r.requiredVersion {
			continue
		}

		runner.EmitIssue(
			r,
			"`policy` on `aws_sqs_queue_policy` must set top-level `Version` to `2012-10-17`; AWS may time out without it.",
			policy.Expr.Range(),
		)
	}

	return nil
}

func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) extractVersion(runner tflint.Runner, expr hcl.Expression) (string, bool, error) {
	if version, ok, err := r.extractVersionFromJSONEncode(runner, expr); ok || err != nil {
		return version, ok, err
	}

	var version string
	resolved := false

	err := runner.EvaluateExpr(expr, func(value cty.Value) error {
		if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
			return nil
		}

		extractedVersion, ok := r.extractVersionFromJSONDocument(value.AsString())
		if !ok {
			return nil
		}

		version = extractedVersion
		resolved = true
		return nil
	}, nil)
	if err != nil {
		return "", false, err
	}

	return version, resolved, nil
}

func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) extractVersionFromJSONEncode(runner tflint.Runner, expr hcl.Expression) (string, bool, error) {
	call, ok := expr.(*hclsyntax.FunctionCallExpr)
	if !ok || call.Name != "jsonencode" || len(call.Args) != 1 {
		return "", false, nil
	}

	object, ok := call.Args[0].(*hclsyntax.ObjectConsExpr)
	if !ok {
		return "", false, nil
	}

	for _, item := range object.Items {
		if hcl.ExprAsKeyword(item.KeyExpr) != "Version" {
			continue
		}

		var version string
		resolved := false

		err := runner.EvaluateExpr(item.ValueExpr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			version = value.AsString()
			resolved = true
			return nil
		}, nil)
		if err != nil {
			return "", false, err
		}
		if !resolved {
			return "", false, nil
		}

		return version, true, nil
	}

	return "", true, nil
}

func (r *AwsSQSQueuePolicyMissingCurrentVersionRule) extractVersionFromJSONDocument(policyDocument string) (string, bool) {
	var document map[string]any
	if err := json.Unmarshal([]byte(policyDocument), &document); err != nil {
		return "", false
	}

	version, ok := document["Version"].(string)
	if !ok {
		return "", true
	}

	return version, true
}
