package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsEFSFileSystemKMSKeyWithoutEncryptedRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "kms key without encrypted attribute",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token = "example"
  kms_key_id     = "arn:aws:kms:us-east-1:123456789012:key/example"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEFSFileSystemKMSKeyWithoutEncryptedRule(),
					Message: "`kms_key_id` can only be set when `encrypted = true` on `aws_efs_file_system`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 20},
						End:      hcl.Pos{Line: 4, Column: 68},
					},
				},
			},
		},
		{
			Name: "kms key with encrypted false",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token = "example"
  encrypted      = false
  kms_key_id     = "arn:aws:kms:us-east-1:123456789012:key/example"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsEFSFileSystemKMSKeyWithoutEncryptedRule(),
					Message: "`kms_key_id` can only be set when `encrypted = true` on `aws_efs_file_system`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 20},
						End:      hcl.Pos{Line: 5, Column: 68},
					},
				},
			},
		},
		{
			Name: "kms key with encrypted true",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token = "example"
  encrypted      = true
  kms_key_id     = "arn:aws:kms:us-east-1:123456789012:key/example"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "encrypted false without kms key",
			Content: `
resource "aws_efs_file_system" "this" {
  creation_token = "example"
  encrypted      = false
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsEFSFileSystemKMSKeyWithoutEncryptedRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
