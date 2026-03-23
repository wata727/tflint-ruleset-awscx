package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "s3 bucket acl expected bucket owner",
			Content: `
resource "aws_s3_bucket_acl" "this" {
  bucket                = "example-bucket"
  acl                   = "private"
  expected_bucket_owner = "123456789012"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule(),
					Message: "`expected_bucket_owner` on `aws_s3_bucket_acl` is deprecated; remove it from configuration.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 27},
						End:      hcl.Pos{Line: 5, Column: 41},
					},
				},
			},
		},
		{
			Name: "s3 bucket logging expected bucket owner",
			Content: `
resource "aws_s3_bucket_logging" "this" {
  bucket                = "example-bucket"
  target_bucket         = "log-bucket"
  target_prefix         = "log/"
  expected_bucket_owner = "123456789012"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule(),
					Message: "`expected_bucket_owner` on `aws_s3_bucket_logging` is deprecated; remove it from configuration.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 27},
						End:      hcl.Pos{Line: 6, Column: 41},
					},
				},
			},
		},
		{
			Name: "s3 bucket versioning without expected bucket owner",
			Content: `
resource "aws_s3_bucket_versioning" "this" {
  bucket = "example-bucket"

  versioning_configuration {
    status = "Enabled"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "plain s3 bucket ignored",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketConfigurationExpectedBucketOwnerDeprecatedRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
