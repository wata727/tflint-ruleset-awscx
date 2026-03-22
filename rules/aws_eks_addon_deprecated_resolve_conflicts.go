package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEKSAddonDeprecatedResolveConflictsRule warns when deprecated resolve_conflicts is used.
type AwsEKSAddonDeprecatedResolveConflictsRule struct {
	tflint.DefaultRule

	resourceType string
	attribute    string
}

// NewAwsEKSAddonDeprecatedResolveConflictsRule returns a new rule.
func NewAwsEKSAddonDeprecatedResolveConflictsRule() *AwsEKSAddonDeprecatedResolveConflictsRule {
	return &AwsEKSAddonDeprecatedResolveConflictsRule{
		resourceType: "aws_eks_addon",
		attribute:    "resolve_conflicts",
	}
}

// Name returns the rule name.
func (r *AwsEKSAddonDeprecatedResolveConflictsRule) Name() string {
	return "awscx_eks_addon_deprecated_resolve_conflicts"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsEKSAddonDeprecatedResolveConflictsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsEKSAddonDeprecatedResolveConflictsRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link.
func (r *AwsEKSAddonDeprecatedResolveConflictsRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_addon"
}

// Check warns when deprecated resolve_conflicts is configured on aws_eks_addon.
func (r *AwsEKSAddonDeprecatedResolveConflictsRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{{Name: r.attribute}},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes[r.attribute]
		if !exists {
			continue
		}

		runner.EmitIssue(
			r,
			"`resolve_conflicts` on `aws_eks_addon` is deprecated; use `resolve_conflicts_on_create` and `resolve_conflicts_on_update` instead.",
			attribute.Expr.Range(),
		)
	}

	return nil
}
