package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedPolicyRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline policy attribute",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = []
  })
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedPolicyRule(),
					Message: "`policy` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_policy` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 12},
						End:      hcl.Pos{Line: 7, Column: 5},
					},
				},
			},
		},
		{
			Name: "standalone bucket policy resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_policy" "this" {
  bucket = aws_s3_bucket.this.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = []
  })
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without policy",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedPolicyRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
