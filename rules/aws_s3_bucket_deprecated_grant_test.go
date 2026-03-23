package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedGrantRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline grant block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"

  grant {
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
    id          = "canonical-user-id"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedGrantRule(),
					Message: "`grant` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_acl` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 8},
					},
				},
			},
		},
		{
			Name: "standalone bucket acl resource",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}

resource "aws_s3_bucket_acl" "this" {
  bucket = aws_s3_bucket.this.id

  access_control_policy {
    owner {
      id = "canonical-user-id"
    }

    grant {
      permission = "FULL_CONTROL"

      grantee {
        id   = "canonical-user-id"
        type = "CanonicalUser"
      }
    }
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without grant block",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedGrantRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
