package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedObjectLockConfigurationRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline object lock configuration block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket              = "example-bucket"
  object_lock_enabled = true

  object_lock_configuration {
    object_lock_enabled = "Enabled"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedObjectLockConfigurationRule(),
					Message: "`object_lock_configuration` on `aws_s3_bucket` is deprecated; use `object_lock_enabled` with `aws_s3_bucket_object_lock_configuration` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 28},
					},
				},
			},
		},
		{
			Name: "standalone object lock configuration resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket              = "example-bucket"
  object_lock_enabled = true
}

resource "aws_s3_bucket_object_lock_configuration" "this" {
  bucket = aws_s3_bucket.this.id

  rule {
    default_retention {
      mode = "GOVERNANCE"
      days = 1
    }
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without object lock configuration block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedObjectLockConfigurationRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
