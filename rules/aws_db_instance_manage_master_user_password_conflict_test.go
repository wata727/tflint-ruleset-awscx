package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceManageMasterUserPasswordConflictRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "managed password with password",
			Content: `
resource "aws_db_instance" "this" {
  identifier                  = "example-db"
  allocated_storage           = 20
  engine                      = "postgres"
  instance_class              = "db.t3.micro"
  username                    = "example"
  manage_master_user_password = true
  password                    = "exampleexample"
  skip_final_snapshot         = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceManageMasterUserPasswordConflictRule(),
					Message: "`password` cannot be set when `manage_master_user_password = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 33},
						End:      hcl.Pos{Line: 9, Column: 49},
					},
				},
			},
		},
		{
			Name: "managed password with write only password",
			Content: `
resource "aws_db_instance" "this" {
  identifier                  = "example-db"
  allocated_storage           = 20
  engine                      = "postgres"
  instance_class              = "db.t3.micro"
  username                    = "example"
  manage_master_user_password = true
  password_wo                 = "exampleexample"
  password_wo_version         = 1
  skip_final_snapshot         = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceManageMasterUserPasswordConflictRule(),
					Message: "`password_wo` cannot be set when `manage_master_user_password = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 33},
						End:      hcl.Pos{Line: 9, Column: 49},
					},
				},
			},
		},
		{
			Name: "managed password without inline password",
			Content: `
resource "aws_db_instance" "this" {
  identifier                  = "example-db"
  allocated_storage           = 20
  engine                      = "postgres"
  instance_class              = "db.t3.micro"
  username                    = "example"
  manage_master_user_password = true
  skip_final_snapshot         = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "manual password without managed password",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 20
  engine              = "postgres"
  instance_class      = "db.t3.micro"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "managed password with both password forms",
			Content: `
resource "aws_db_instance" "this" {
  identifier                  = "example-db"
  allocated_storage           = 20
  engine                      = "postgres"
  instance_class              = "db.t3.micro"
  username                    = "example"
  manage_master_user_password = true
  password                    = "exampleexample"
  password_wo                 = "anotherexample"
  password_wo_version         = 1
  skip_final_snapshot         = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceManageMasterUserPasswordConflictRule(),
					Message: "`password` cannot be set when `manage_master_user_password = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 9, Column: 33},
						End:      hcl.Pos{Line: 9, Column: 49},
					},
				},
				{
					Rule:    NewAwsDBInstanceManageMasterUserPasswordConflictRule(),
					Message: "`password_wo` cannot be set when `manage_master_user_password = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 33},
						End:      hcl.Pos{Line: 10, Column: 49},
					},
				},
			},
		},
	}

	rule := NewAwsDBInstanceManageMasterUserPasswordConflictRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
