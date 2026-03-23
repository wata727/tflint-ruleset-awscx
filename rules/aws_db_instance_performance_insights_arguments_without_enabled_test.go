package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "kms key without performance insights enabled",
			Content: `
resource "aws_db_instance" "this" {
  identifier                        = "example-db"
  allocated_storage                 = 20
  engine                            = "postgres"
  instance_class                    = "db.t3.micro"
  username                          = "example"
  password                          = "exampleexample"
  skip_final_snapshot               = true
  performance_insights_kms_key_id   = aws_kms_key.example.arn
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule(),
					Message: "`performance_insights_kms_key_id` cannot be set unless `performance_insights_enabled = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 39},
						End:      hcl.Pos{Line: 10, Column: 62},
					},
				},
			},
		},
		{
			Name: "retention period with performance insights disabled",
			Content: `
resource "aws_db_instance" "this" {
  identifier                            = "example-db"
  allocated_storage                     = 20
  engine                                = "postgres"
  instance_class                        = "db.t3.micro"
  username                              = "example"
  password                              = "exampleexample"
  skip_final_snapshot                   = true
  performance_insights_enabled          = false
  performance_insights_retention_period = 93
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule(),
					Message: "`performance_insights_retention_period` cannot be set unless `performance_insights_enabled = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 43},
						End:      hcl.Pos{Line: 11, Column: 45},
					},
				},
			},
		},
		{
			Name: "both dependent arguments with performance insights disabled",
			Content: `
resource "aws_db_instance" "this" {
  identifier                            = "example-db"
  allocated_storage                     = 20
  engine                                = "postgres"
  instance_class                        = "db.t3.micro"
  username                              = "example"
  password                              = "exampleexample"
  skip_final_snapshot                   = true
  performance_insights_enabled          = false
  performance_insights_kms_key_id       = aws_kms_key.example.arn
  performance_insights_retention_period = 93
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule(),
					Message: "`performance_insights_kms_key_id` cannot be set unless `performance_insights_enabled = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 43},
						End:      hcl.Pos{Line: 11, Column: 66},
					},
				},
				{
					Rule:    NewAwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule(),
					Message: "`performance_insights_retention_period` cannot be set unless `performance_insights_enabled = true` on `aws_db_instance`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 12, Column: 43},
						End:      hcl.Pos{Line: 12, Column: 45},
					},
				},
			},
		},
		{
			Name: "performance insights enabled with dependent arguments",
			Content: `
resource "aws_db_instance" "this" {
  identifier                            = "example-db"
  allocated_storage                     = 20
  engine                                = "postgres"
  instance_class                        = "db.t3.micro"
  username                              = "example"
  password                              = "exampleexample"
  skip_final_snapshot                   = true
  performance_insights_enabled          = true
  performance_insights_kms_key_id       = aws_kms_key.example.arn
  performance_insights_retention_period = 93
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "performance insights disabled without dependent arguments",
			Content: `
resource "aws_db_instance" "this" {
  identifier                   = "example-db"
  allocated_storage            = 20
  engine                       = "postgres"
  instance_class               = "db.t3.micro"
  username                     = "example"
  password                     = "exampleexample"
  skip_final_snapshot          = true
  performance_insights_enabled = false
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstancePerformanceInsightsArgumentsWithoutEnabledRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
