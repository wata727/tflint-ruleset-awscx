package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceMissingIOPSRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "gp3 without iops",
			Content: `
resource "aws_db_instance" "this" {
  identifier         = "example-db"
  allocated_storage  = 100
  engine             = "postgres"
  instance_class     = "db.t3.medium"
  username           = "example"
  password           = "exampleexample"
  skip_final_snapshot = true
  storage_type       = "gp3"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceMissingIOPSRule(),
					Message: "`iops` must be set when `storage_type` is `io1`, `io2`, or `gp3`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 24},
						End:      hcl.Pos{Line: 10, Column: 29},
					},
				},
			},
		},
		{
			Name: "io1 with iops",
			Content: `
resource "aws_db_instance" "this" {
  identifier          = "example-db"
  allocated_storage   = 100
  engine              = "postgres"
  instance_class      = "db.t3.medium"
  username            = "example"
  password            = "exampleexample"
  skip_final_snapshot = true
  storage_type        = "io1"
  iops                = 1000
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "gp2 without iops",
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
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceMissingIOPSRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
