package rules

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// AwsLambdaFunctionZipRequiredAttributesRule checks Zip-based Lambda functions declare required attributes.
type AwsLambdaFunctionZipRequiredAttributesRule struct {
	tflint.DefaultRule

	resourceType    string
	packageTypeAttr string
	handlerAttr     string
	runtimeAttr     string
	expectedPkgType string
}

// NewAwsLambdaFunctionZipRequiredAttributesRule returns a new rule.
func NewAwsLambdaFunctionZipRequiredAttributesRule() *AwsLambdaFunctionZipRequiredAttributesRule {
	return &AwsLambdaFunctionZipRequiredAttributesRule{
		resourceType:    "aws_lambda_function",
		packageTypeAttr: "package_type",
		handlerAttr:     "handler",
		runtimeAttr:     "runtime",
		expectedPkgType: "ZIP",
	}
}

// Name returns the rule name.
func (r *AwsLambdaFunctionZipRequiredAttributesRule) Name() string {
	return "awscx_lambda_function_zip_required_attributes"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLambdaFunctionZipRequiredAttributesRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLambdaFunctionZipRequiredAttributesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLambdaFunctionZipRequiredAttributesRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function"
}

// Check reports Zip Lambda functions that omit handler or runtime.
func (r *AwsLambdaFunctionZipRequiredAttributesRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.packageTypeAttr},
			{Name: r.handlerAttr},
			{Name: r.runtimeAttr},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		packageTypeAttr, hasPackageType := resource.Body.Attributes[r.packageTypeAttr]
		_, hasHandler := resource.Body.Attributes[r.handlerAttr]
		_, hasRuntime := resource.Body.Attributes[r.runtimeAttr]

		if hasHandler && hasRuntime {
			continue
		}

		if !hasPackageType {
			r.emitIssues(runner, resource.DefRange, hasHandler, hasRuntime)
			continue
		}

		err := runner.EvaluateExpr(packageTypeAttr.Expr, func(value cty.Value) error {
			if !value.IsKnown() || value.IsNull() || !value.Type().Equals(cty.String) {
				return nil
			}

			if strings.ToUpper(strings.TrimSpace(value.AsString())) != r.expectedPkgType {
				return nil
			}

			r.emitIssues(runner, packageTypeAttr.Expr.Range(), hasHandler, hasRuntime)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *AwsLambdaFunctionZipRequiredAttributesRule) emitIssues(runner tflint.Runner, issueRange hcl.Range, hasHandler bool, hasRuntime bool) {
	if !hasHandler {
		runner.EmitIssue(
			r,
			"`handler` must be set when `package_type` is `Zip` on `aws_lambda_function`.",
			issueRange,
		)
	}

	if !hasRuntime {
		runner.EmitIssue(
			r,
			"`runtime` must be set when `package_type` is `Zip` on `aws_lambda_function`.",
			issueRange,
		)
	}
}
