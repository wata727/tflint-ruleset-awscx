package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/wata727/tflint-ruleset-awscx/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "awscx",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewAwsDBInstanceMissingIOPSRule(),
				rules.NewAwsDynamoDBTableInvalidStreamViewTypeRule(),
				rules.NewAwsEBSVolumeMissingIOPSRule(),
				rules.NewAwsEFSFileSystemMissingProvisionedThroughputRule(),
				rules.NewAwsInstanceIMDSv2OptionalTokensRule(),
				rules.NewAwsLaunchTemplateIMDSv2OptionalTokensRule(),
				rules.NewAwsS3BucketDeprecatedACLRule(),
				rules.NewAwsS3BucketDeprecatedServerSideEncryptionConfigurationRule(),
				rules.NewAwsS3BucketDeprecatedVersioningRule(),
				rules.NewAwsSecurityGroupInvalidProtocolRule(),
			},
		},
	})
}
