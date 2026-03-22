package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedVersioningRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline versioning block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"

  versioning {
    enabled = true
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedVersioningRule(),
					Message: "`versioning` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_versioning` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 13},
					},
				},
			},
		},
		{
			Name: "standalone versioning resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_versioning" "this" {
  bucket = aws_s3_bucket.this.id

  versioning_configuration {
    status = "Enabled"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without versioning block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedVersioningRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
