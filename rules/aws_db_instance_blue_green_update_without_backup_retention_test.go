package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "blue green update without backup retention period",
			Content: `
resource "aws_db_instance" "this" {
  identifier         = "example-db"
  allocated_storage  = 20
  engine             = "mysql"
  instance_class     = "db.t3.micro"
  username           = "example"
  password           = "exampleexample"
  skip_final_snapshot = true

  blue_green_update {
    enabled = true
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule(),
					Message: "`blue_green_update.enabled = true` requires `backup_retention_period` greater than `0` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 15},
						End:      hcl.Pos{Line: 12, Column: 19},
					},
				},
			},
		},
		{
			Name: "blue green update with zero backup retention period",
			Content: `
resource "aws_db_instance" "this" {
  identifier              = "example-db"
  allocated_storage       = 20
  engine                  = "mysql"
  instance_class          = "db.t3.micro"
  username                = "example"
  password                = "exampleexample"
  skip_final_snapshot     = true
  backup_retention_period = 0

  blue_green_update {
    enabled = true
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule(),
					Message: "`blue_green_update.enabled = true` requires `backup_retention_period` greater than `0` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 29},
						End:      hcl.Pos{Line: 10, Column: 30},
					},
				},
			},
		},
		{
			Name: "blue green update with positive backup retention period",
			Content: `
resource "aws_db_instance" "this" {
  identifier              = "example-db"
  allocated_storage       = 20
  engine                  = "mysql"
  instance_class          = "db.t3.micro"
  username                = "example"
  password                = "exampleexample"
  skip_final_snapshot     = true
  backup_retention_period = 7

  blue_green_update {
    enabled = true
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "disabled blue green update does not require backup retention",
			Content: `
resource "aws_db_instance" "this" {
  identifier              = "example-db"
  allocated_storage       = 20
  engine                  = "mysql"
  instance_class          = "db.t3.micro"
  username                = "example"
  password                = "exampleexample"
  skip_final_snapshot     = true
  backup_retention_period = 0

  blue_green_update {
    enabled = false
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown backup retention is skipped",
			Content: `
variable "backup_retention_period" {
  type = number
}

resource "aws_db_instance" "this" {
  identifier              = "example-db"
  allocated_storage       = 20
  engine                  = "mysql"
  instance_class          = "db.t3.micro"
  username                = "example"
  password                = "exampleexample"
  skip_final_snapshot     = true
  backup_retention_period = var.backup_retention_period

  blue_green_update {
    enabled = true
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceBlueGreenUpdateWithoutBackupRetentionRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
