package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedWebsiteRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline website block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedWebsiteRule(),
					Message: "`website` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_website_configuration` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 10},
					},
				},
			},
		},
		{
			Name: "standalone website configuration resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_website_configuration" "this" {
  bucket = aws_s3_bucket.this.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without website block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedWebsiteRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
