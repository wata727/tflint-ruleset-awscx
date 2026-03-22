package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "deprecated elastic gpu block",
			Content: `
resource "aws_launch_template" "this" {
  name_prefix = "example"

  elastic_gpu_specifications {
    type = "eg1.medium"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule(),
					Message: "`elastic_gpu_specifications` on `aws_launch_template` is deprecated because Amazon Elastic Graphics reached end of life on January 8, 2024.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 29},
					},
				},
			},
		},
		{
			Name: "without elastic gpu block",
			Content: `
resource "aws_launch_template" "this" {
  name_prefix = "example"
  image_id    = "ami-12345678"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "other resource with similarly named block",
			Content: `
resource "aws_instance" "this" {
  ami           = "ami-12345678"
  instance_type = "t3.micro"

  elastic_gpu_specifications {
    type = "eg1.medium"
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
