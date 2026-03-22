package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedACLRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "deprecated inline acl",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
  acl    = "private"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedACLRule(),
					Message: "`acl` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_acl` and `aws_s3_bucket_ownership_controls` when ACLs are required.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 12},
						End:      hcl.Pos{Line: 4, Column: 21},
					},
				},
			},
		},
		{
			Name: "no inline acl",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedACLRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
