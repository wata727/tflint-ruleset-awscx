package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEKSAddonDeprecatedResolveConflictsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "deprecated resolve_conflicts attribute",
			Content: `
resource "aws_eks_addon" "this" {
  cluster_name      = "example"
  addon_name        = "vpc-cni"
  resolve_conflicts = "OVERWRITE"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEKSAddonDeprecatedResolveConflictsRule(),
					Message: "`resolve_conflicts` on `aws_eks_addon` is deprecated; use `resolve_conflicts_on_create` and `resolve_conflicts_on_update` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 23},
						End:      hcl.Pos{Line: 5, Column: 34},
					},
				},
			},
		},
		{
			Name: "replacement attributes only",
			Content: `
resource "aws_eks_addon" "this" {
  cluster_name                = "example"
  addon_name                  = "vpc-cni"
  resolve_conflicts_on_create = "OVERWRITE"
  resolve_conflicts_on_update = "PRESERVE"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no conflict attributes",
			Content: `
resource "aws_eks_addon" "this" {
  cluster_name = "example"
  addon_name   = "vpc-cni"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEKSAddonDeprecatedResolveConflictsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
