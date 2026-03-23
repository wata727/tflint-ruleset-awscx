package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceEnhancedMonitoringRoleRequirementsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "enhanced monitoring interval without role",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 20
  engine              = "postgres"
  instance_class      = "db.t3.micro"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  monitoring_interval = 60
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceEnhancedMonitoringRoleRequirementsRule(),
					Message: "`monitoring_role_arn` must be set when `monitoring_interval` enables Enhanced Monitoring on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 25},
						End:      hcl.Pos{Line: 10, Column: 27},
					},
				},
			},
		},
		{
			Name: "monitoring role with interval zero",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 20
  engine              = "postgres"
  instance_class      = "db.t3.micro"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  monitoring_interval = 0
  monitoring_role_arn = aws_iam_role.rds_monitoring.arn
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceEnhancedMonitoringRoleRequirementsRule(),
					Message: "`monitoring_role_arn` cannot be set when `monitoring_interval = 0` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 25},
						End:      hcl.Pos{Line: 11, Column: 56},
					},
				},
			},
		},
		{
			Name: "enhanced monitoring interval with role",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 20
  engine              = "postgres"
  instance_class      = "db.t3.micro"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  monitoring_interval = 60
  monitoring_role_arn = aws_iam_role.rds_monitoring.arn
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "interval zero without role",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 20
  engine              = "postgres"
  instance_class      = "db.t3.micro"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  monitoring_interval = 0
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "role without explicit interval is skipped conservatively",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 20
  engine              = "postgres"
  instance_class      = "db.t3.micro"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  monitoring_role_arn = aws_iam_role.rds_monitoring.arn
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceEnhancedMonitoringRoleRequirementsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
