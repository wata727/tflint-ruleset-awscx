package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsDBInstanceManageMasterUserPasswordConflictRule checks invalid RDS master password management combinations.
type AwsDBInstanceManageMasterUserPasswordConflictRule struct {
	tflint.DefaultRule

	resourceType               string
	managedPasswordAttribute   string
	passwordAttribute          string
	writeOnlyPasswordAttribute string
}

// NewAwsDBInstanceManageMasterUserPasswordConflictRule returns a new rule.
func NewAwsDBInstanceManageMasterUserPasswordConflictRule() *AwsDBInstanceManageMasterUserPasswordConflictRule {
	return &AwsDBInstanceManageMasterUserPasswordConflictRule{
		resourceType:               "aws_db_instance",
		managedPasswordAttribute:   "manage_master_user_password",
		passwordAttribute:          "password",
		writeOnlyPasswordAttribute: "password_wo",
	}
}

// Name returns the rule name.
func (r *AwsDBInstanceManageMasterUserPasswordConflictRule) Name() string {
	return "awscx_db_instance_manage_master_user_password_conflict"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsDBInstanceManageMasterUserPasswordConflictRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsDBInstanceManageMasterUserPasswordConflictRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsDBInstanceManageMasterUserPasswordConflictRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance"
}

// Check reports conflicting password arguments when RDS-managed master credentials are enabled.
func (r *AwsDBInstanceManageMasterUserPasswordConflictRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.managedPasswordAttribute},
			{Name: r.passwordAttribute},
			{Name: r.writeOnlyPasswordAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		managedPassword, hasManagedPassword := resource.Body.Attributes[r.managedPasswordAttribute]
		if !hasManagedPassword {
			continue
		}

		password, hasPassword := resource.Body.Attributes[r.passwordAttribute]
		writeOnlyPassword, hasWriteOnlyPassword := resource.Body.Attributes[r.writeOnlyPasswordAttribute]

		err := runner.EvaluateExpr(managedPassword.Expr, func(enabled bool) error {
			if !enabled {
				return nil
			}

			if hasPassword {
				runner.EmitIssue(
					r,
					"`password` cannot be set when `manage_master_user_password = true` on `aws_db_instance`.",
					password.Expr.Range(),
				)
			}

			if hasWriteOnlyPassword {
				runner.EmitIssue(
					r,
					"`password_wo` cannot be set when `manage_master_user_password = true` on `aws_db_instance`.",
					writeOnlyPassword.Expr.Range(),
				)
			}

			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
