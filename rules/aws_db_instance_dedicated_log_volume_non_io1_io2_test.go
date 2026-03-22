package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "dedicated log volume without storage type",
			Content: `
resource "aws_db_instance" "this" {
  identifier           = "example-db"
  allocated_storage    = 100
  engine               = "postgres"
  instance_class       = "db.t3.medium"
  username             = "example"
  password             = "exampleexample"
  skip_final_snapshot  = true
  dedicated_log_volume = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule(),
					Message: "`dedicated_log_volume` can only be set when `storage_type` is `io1` or `io2`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 26},
						End:      hcl.Pos{Line: 10, Column: 30},
					},
				},
			},
		},
		{
			Name: "dedicated log volume with gp3",
			Content: `
resource "aws_db_instance" "this" {
  identifier           = "example-db"
  allocated_storage    = 100
  engine               = "postgres"
  instance_class       = "db.t3.medium"
  username             = "example"
  password             = "exampleexample"
  skip_final_snapshot  = true
  storage_type         = "gp3"
  dedicated_log_volume = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule(),
					Message: "`dedicated_log_volume` can only be set when `storage_type` is `io1` or `io2`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 11, Column: 26},
						End:      hcl.Pos{Line: 11, Column: 30},
					},
				},
			},
		},
		{
			Name: "dedicated log volume with io2",
			Content: `
resource "aws_db_instance" "this" {
  identifier           = "example-db"
  allocated_storage    = 100
  engine               = "postgres"
  instance_class       = "db.t3.medium"
  username             = "example"
  password             = "exampleexample"
  skip_final_snapshot  = true
  storage_type         = "io2"
  iops                 = 1000
  dedicated_log_volume = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "dedicated log volume disabled with gp3",
			Content: `
resource "aws_db_instance" "this" {
  identifier           = "example-db"
  allocated_storage    = 100
  engine               = "postgres"
  instance_class       = "db.t3.medium"
  username             = "example"
  password             = "exampleexample"
  skip_final_snapshot  = true
  storage_type         = "gp3"
  dedicated_log_volume = false
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsDBInstanceDedicatedLogVolumeNonIO1IO2Rule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
