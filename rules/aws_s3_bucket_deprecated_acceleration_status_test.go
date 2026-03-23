package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedAccelerationStatusRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline acceleration_status attribute",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket              = "example-bucket"
  acceleration_status = "Enabled"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedAccelerationStatusRule(),
					Message: "`acceleration_status` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_accelerate_configuration` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 25},
						End:      hcl.Pos{Line: 4, Column: 34},
					},
				},
			},
		},
		{
			Name: "standalone accelerate configuration resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_accelerate_configuration" "this" {
  bucket = aws_s3_bucket.this.id
  status = "Enabled"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without acceleration_status",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedAccelerationStatusRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
