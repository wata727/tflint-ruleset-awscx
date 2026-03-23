package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceDatabaseInsightsAdvancedRequirementsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "advanced mode without performance insights settings",
			Content: `
resource "aws_db_instance" "this" {
  identifier             = "example-db"
  allocated_storage      = 20
  engine                 = "postgres"
  instance_class         = "db.t3.micro"
  username               = "example"
  password               = "exampleexample"
  skip_final_snapshot    = true
  database_insights_mode = "advanced"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceDatabaseInsightsAdvancedRequirementsRule(),
					Message: "`database_insights_mode = \"advanced\"` requires `performance_insights_enabled = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 28},
						End:      hcl.Pos{Line: 10, Column: 38},
					},
				},
				{
					Rule:    NewAwsDBInstanceDatabaseInsightsAdvancedRequirementsRule(),
					Message: "`database_insights_mode = \"advanced\"` requires `performance_insights_retention_period` to be at least `465` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 28},
						End:      hcl.Pos{Line: 10, Column: 38},
					},
				},
			},
		},
		{
			Name: "advanced mode with performance insights disabled and short retention",
			Content: `
resource "aws_db_instance" "this" {
  identifier                            = "example-db"
  allocated_storage                     = 20
  engine                                = "postgres"
  instance_class                        = "db.t3.micro"
  username                              = "example"
  password                              = "exampleexample"
  skip_final_snapshot                   = true
  database_insights_mode                = "advanced"
  performance_insights_enabled          = false
  performance_insights_retention_period = 93
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceDatabaseInsightsAdvancedRequirementsRule(),
					Message: "`database_insights_mode = \"advanced\"` requires `performance_insights_enabled = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 43},
						End:      hcl.Pos{Line: 11, Column: 48},
					},
				},
				{
					Rule:    NewAwsDBInstanceDatabaseInsightsAdvancedRequirementsRule(),
					Message: "`database_insights_mode = \"advanced\"` requires `performance_insights_retention_period` to be at least `465` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 43},
						End:      hcl.Pos{Line: 12, Column: 45},
					},
				},
			},
		},
		{
			Name: "advanced mode with enabled performance insights and valid retention",
			Content: `
resource "aws_db_instance" "this" {
  identifier                            = "example-db"
  allocated_storage                     = 20
  engine                                = "postgres"
  instance_class                        = "db.t3.micro"
  username                              = "example"
  password                              = "exampleexample"
  skip_final_snapshot                   = true
  database_insights_mode                = "advanced"
  performance_insights_enabled          = true
  performance_insights_retention_period = 465
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "standard mode does not require advanced settings",
			Content: `
resource "aws_db_instance" "this" {
  identifier                            = "example-db"
  allocated_storage                     = 20
  engine                                = "postgres"
  instance_class                        = "db.t3.micro"
  username                              = "example"
  password                              = "exampleexample"
  skip_final_snapshot                   = true
  database_insights_mode                = "standard"
  performance_insights_enabled          = false
  performance_insights_retention_period = 7
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "advanced mode with unknown retention is skipped",
			Content: `
variable "retention_days" {
  type = number
}

resource "aws_db_instance" "this" {
  identifier                            = "example-db"
  allocated_storage                     = 20
  engine                                = "postgres"
  instance_class                        = "db.t3.micro"
  username                              = "example"
  password                              = "exampleexample"
  skip_final_snapshot                   = true
  database_insights_mode                = "advanced"
  performance_insights_enabled          = true
  performance_insights_retention_period = var.retention_days
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceDatabaseInsightsAdvancedRequirementsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
