package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketDeprecatedReplicationConfigurationRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "inline replication configuration block",
			Content: `
resource "aws_s3_bucket" "source" {
  bucket = "example-source-bucket"

  replication_configuration {
    role = aws_iam_role.replication.arn

    rules {
      id     = "all"
      status = "Enabled"

      destination {
        bucket = aws_s3_bucket.destination.arn
      }
    }
  }
}

resource "aws_s3_bucket" "destination" {
  bucket = "example-destination-bucket"
}

resource "aws_iam_role" "replication" {
  name = "example-replication-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "s3.amazonaws.com"
      }
    }]
  })
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsS3BucketDeprecatedReplicationConfigurationRule(),
					Message: "`replication_configuration` on `aws_s3_bucket` is deprecated; use `aws_s3_bucket_replication_configuration` instead.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 28},
					},
				},
			},
		},
		{
			Name: "standalone replication configuration resource",
			Content: `
resource "aws_s3_bucket" "source" {
  bucket = "example-source-bucket"
}

resource "aws_s3_bucket" "destination" {
  bucket = "example-destination-bucket"
}

resource "aws_s3_bucket_replication_configuration" "this" {
  bucket = aws_s3_bucket.source.id
  role   = aws_iam_role.replication.arn

  rule {
    id     = "all"
    status = "Enabled"

    destination {
      bucket = aws_s3_bucket.destination.arn
    }
  }
}

resource "aws_iam_role" "replication" {
  name = "example-replication-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "s3.amazonaws.com"
      }
    }]
  })
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket without replication configuration",
			Content: `
resource "aws_s3_bucket" "this" {
  bucket = "example-bucket"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsS3BucketDeprecatedReplicationConfigurationRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
