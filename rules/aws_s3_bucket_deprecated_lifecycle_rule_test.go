package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedLifecycleRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline lifecycle_rule block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"

  lifecycle_rule {
    id      = "expire-logs"
    enabled = true

    expiration {
      days = 30
    }
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedLifecycleRule(),
					Message: "`lifecycle_rule` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_lifecycle_configuration` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 17},
					},
				},
			},
		},
		{
			Name: "standalone lifecycle configuration resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_lifecycle_configuration" "this" {
  bucket = aws_s3_bucket.this.id

  rule {
    id     = "expire-logs"
    status = "Enabled"

    expiration {
      days = 30
    }
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without lifecycle_rule block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedLifecycleRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
