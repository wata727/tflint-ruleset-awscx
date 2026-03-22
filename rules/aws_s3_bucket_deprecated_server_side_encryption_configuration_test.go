package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedServerSideEncryptionConfigurationRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline encryption block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedServerSideEncryptionConfigurationRule(),
					Message: "`server_side_encryption_configuration` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_server_side_encryption_configuration` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 39},
					},
				},
			},
		},
		{
			Name: "standalone encryption resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  bucket = aws_s3_bucket.this.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without encryption block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedServerSideEncryptionConfigurationRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
