package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceStorageThroughputNonGP3Rule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "storage throughput without storage type",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 100
  engine              = "postgres"
  instance_class      = "db.t3.medium"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  storage_throughput  = 250
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceStorageThroughputNonGP3Rule(),
					Message: "`storage_throughput` can only be set when `storage_type = \"gp3\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 25},
						End:      hcl.Pos{Line: 10, Column: 28},
					},
				},
			},
		},
		{
			Name: "storage throughput with gp2",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 100
  engine              = "postgres"
  instance_class      = "db.t3.medium"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  storage_type        = "gp2"
  storage_throughput  = 250
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceStorageThroughputNonGP3Rule(),
					Message: "`storage_throughput` can only be set when `storage_type = \"gp3\"`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 25},
						End:      hcl.Pos{Line: 11, Column: 28},
					},
				},
			},
		},
		{
			Name: "storage throughput with gp3",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 400
  engine              = "postgres"
  instance_class      = "db.t3.medium"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  storage_type        = "gp3"
  iops                = 12000
  storage_throughput  = 500
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceStorageThroughputNonGP3Rule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
