package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedRequestPayerRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline request payer attribute",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket        = "example-bucket"
  request_payer = "Requester"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedRequestPayerRule(),
					Message: "`request_payer` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_request_payment_configuration` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 19},
						End:      hcl.Pos{Line: 4, Column: 30},
					},
				},
			},
		},
		{
			Name: "standalone request payment configuration resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_request_payment_configuration" "this" {
  bucket = aws_s3_bucket.this.id
  payer  = "Requester"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without request payer",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedRequestPayerRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
