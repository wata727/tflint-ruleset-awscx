package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsEFSFileSystemKMSKeyWithoutEncryptedRule checks that kms_key_id is only used with encryption enabled.
type AwsEFSFileSystemKMSKeyWithoutEncryptedRule struct {
	tflint.DefaultRule

	resourceType       string
	encryptedAttribute string
	kmsKeyIDAttribute  string
}

// NewAwsEFSFileSystemKMSKeyWithoutEncryptedRule returns a new rule.
func NewAwsEFSFileSystemKMSKeyWithoutEncryptedRule() *AwsEFSFileSystemKMSKeyWithoutEncryptedRule {
	return &AwsEFSFileSystemKMSKeyWithoutEncryptedRule{
		resourceType:       "aws_efs_file_system",
		encryptedAttribute: "encrypted",
		kmsKeyIDAttribute:  "kms_key_id",
	}
}

// Name returns the rule name.
func (r *AwsEFSFileSystemKMSKeyWithoutEncryptedRule) Name() string {
	return "awscx_efs_file_system_kms_key_without_encrypted"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsEFSFileSystemKMSKeyWithoutEncryptedRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsEFSFileSystemKMSKeyWithoutEncryptedRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsEFSFileSystemKMSKeyWithoutEncryptedRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system"
}

// Check reports when kms_key_id is configured without encryption enabled.
func (r *AwsEFSFileSystemKMSKeyWithoutEncryptedRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: r.encryptedAttribute},
			{Name: r.kmsKeyIDAttribute},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		kmsKeyID, hasKMSKeyID := resource.Body.Attributes[r.kmsKeyIDAttribute]
		if !hasKMSKeyID {
			continue
		}

		encrypted, hasEncrypted := resource.Body.Attributes[r.encryptedAttribute]
		if !hasEncrypted {
			runner.EmitIssue(
				r,
				"`kms_key_id` can only be set when `encrypted = true` on `aws_efs_file_system`.",
				kmsKeyID.Expr.Range(),
			)
			continue
		}

		err := runner.EvaluateExpr(encrypted.Expr, func(value bool) error {
			if value {
				return nil
			}

			runner.EmitIssue(
				r,
				"`kms_key_id` can only be set when `encrypted = true` on `aws_efs_file_system`.",
				kmsKeyID.Expr.Range(),
			)
			return nil
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
