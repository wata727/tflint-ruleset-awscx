package rules

import (
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsLBListenerMutualAuthenticationVerifyRequirementsRule checks aws_lb_listener mutual authentication settings.
type AwsLBListenerMutualAuthenticationVerifyRequirementsRule struct {
	tflint.DefaultRule

	resourceType                  string
	blockType                     string
	modeAttribute                 string
	trustStoreARNAttribute        string
	advertiseTrustStoreCANames    string
	ignoreClientCertificateExpiry string
}

// NewAwsLBListenerMutualAuthenticationVerifyRequirementsRule returns a new rule.
func NewAwsLBListenerMutualAuthenticationVerifyRequirementsRule() *AwsLBListenerMutualAuthenticationVerifyRequirementsRule {
	return &AwsLBListenerMutualAuthenticationVerifyRequirementsRule{
		resourceType:                  "aws_lb_listener",
		blockType:                     "mutual_authentication",
		modeAttribute:                 "mode",
		trustStoreARNAttribute:        "trust_store_arn",
		advertiseTrustStoreCANames:    "advertise_trust_store_ca_names",
		ignoreClientCertificateExpiry: "ignore_client_certificate_expiry",
	}
}

// Name returns the rule name.
func (r *AwsLBListenerMutualAuthenticationVerifyRequirementsRule) Name() string {
	return "awscx_lb_listener_mutual_authentication_verify_requirements"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsLBListenerMutualAuthenticationVerifyRequirementsRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsLBListenerMutualAuthenticationVerifyRequirementsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsLBListenerMutualAuthenticationVerifyRequirementsRule) Link() string {
	return "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener"
}

// Check reports invalid mutual authentication combinations on aws_lb_listener.
func (r *AwsLBListenerMutualAuthenticationVerifyRequirementsRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent(r.resourceType, &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: r.blockType,
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: r.modeAttribute},
						{Name: r.trustStoreARNAttribute},
						{Name: r.advertiseTrustStoreCANames},
						{Name: r.ignoreClientCertificateExpiry},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			mode, exists := block.Body.Attributes[r.modeAttribute]
			if !exists {
				continue
			}

			trustStoreARN, hasTrustStoreARN := block.Body.Attributes[r.trustStoreARNAttribute]
			advertiseTrustStoreCANames, hasAdvertiseTrustStoreCANames := block.Body.Attributes[r.advertiseTrustStoreCANames]
			ignoreClientCertificateExpiry, hasIgnoreClientCertificateExpiry := block.Body.Attributes[r.ignoreClientCertificateExpiry]

			err := runner.EvaluateExpr(mode.Expr, func(value string) error {
				switch strings.ToLower(strings.TrimSpace(value)) {
				case "verify":
					if hasTrustStoreARN {
						return nil
					}

					runner.EmitIssue(
						r,
						"`mutual_authentication.trust_store_arn` must be set when `mutual_authentication.mode` is `\"verify\"`.",
						mode.Expr.Range(),
					)
				case "off", "passthrough":
					if hasTrustStoreARN {
						runner.EmitIssue(
							r,
							"`mutual_authentication.trust_store_arn` is only valid when `mutual_authentication.mode` is `\"verify\"`.",
							trustStoreARN.Expr.Range(),
						)
					}

					if hasAdvertiseTrustStoreCANames {
						runner.EmitIssue(
							r,
							"`mutual_authentication.advertise_trust_store_ca_names` is only valid when `mutual_authentication.mode` is `\"verify\"`.",
							advertiseTrustStoreCANames.Expr.Range(),
						)
					}

					if hasIgnoreClientCertificateExpiry {
						runner.EmitIssue(
							r,
							"`mutual_authentication.ignore_client_certificate_expiry` is only valid when `mutual_authentication.mode` is `\"verify\"`.",
							ignoreClientCertificateExpiry.Expr.Range(),
						)
					}
				}

				return nil
			}, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
