package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsLambdaFunctionZipRequiredAttributesRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "explicit zip missing handler and runtime",
			Content: `
resource "aws_lambda_function" "this" {
  function_name = "example"
  role          = aws_iam_role.this.arn
  package_type  = "Zip"
  filename      = "function.zip"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionZipRequiredAttributesRule(),
					Message: "`handler` must be set when `package_type` is `Zip` on `aws_lambda_function`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 19},
						End:      hcl.Pos{Line: 5, Column: 24},
					},
				},
				{
					Rule:    NewAwsLambdaFunctionZipRequiredAttributesRule(),
					Message: "`runtime` must be set when `package_type` is `Zip` on `aws_lambda_function`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 5, Column: 19},
						End:      hcl.Pos{Line: 5, Column: 24},
					},
				},
			},
		},
		{
			Name: "default zip missing runtime",
			Content: `
resource "aws_lambda_function" "this" {
  function_name = "example"
  role          = aws_iam_role.this.arn
  filename      = "function.zip"
  handler       = "index.handler"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsLambdaFunctionZipRequiredAttributesRule(),
					Message: "`runtime` must be set when `package_type` is `Zip` on `aws_lambda_function`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 38},
					},
				},
			},
		},
		{
			Name: "image package omits handler and runtime",
			Content: `
resource "aws_lambda_function" "this" {
  function_name = "example"
  role          = aws_iam_role.this.arn
  package_type  = "Image"
  image_uri     = "123456789012.dkr.ecr.us-east-1.amazonaws.com/example:latest"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "zip package with required attributes",
			Content: `
resource "aws_lambda_function" "this" {
  function_name = "example"
  role          = aws_iam_role.this.arn
  package_type  = "Zip"
  filename      = "function.zip"
  handler       = "index.handler"
  runtime       = "python3.13"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown package type expression skipped",
			Content: `
variable "package_type" {
  type = string
}

resource "aws_lambda_function" "this" {
  function_name = "example"
  role          = aws_iam_role.this.arn
  package_type  = var.package_type
  filename      = "function.zip"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsLambdaFunctionZipRequiredAttributesRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
