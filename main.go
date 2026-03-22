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
				rules.NewAwsAPIGatewayDeploymentDeprecatedStageManagementRule(),
				rules.NewAwsDBInstanceMissingIOPSRule(),
				rules.NewAwsDBInstanceStorageThroughputNonGP3Rule(),
				rules.NewAwsDynamoDBTableInvalidStreamViewTypeRule(),
				rules.NewAwsEBSVolumeMissingIOPSRule(),
				rules.NewAwsEBSVolumeThroughputNonGP3Rule(),
				rules.NewAwsEKSAddonDeprecatedResolveConflictsRule(),
				rules.NewAwsEFSFileSystemMissingProvisionedThroughputRule(),
				rules.NewAwsInstanceDeprecatedNetworkInterfaceRule(),
				rules.NewAwsLaunchTemplateDeprecatedElasticGPUSpecificationsRule(),
				rules.NewAwsInstanceIMDSv2OptionalTokensRule(),
				rules.NewAwsLaunchTemplateIMDSv2OptionalTokensRule(),
				rules.NewAwsS3BucketDeprecatedACLRule(),
				rules.NewAwsS3BucketDeprecatedLifecycleRule(),
				rules.NewAwsS3BucketDeprecatedLoggingRule(),
				rules.NewAwsS3BucketDeprecatedReplicationConfigurationRule(),
				rules.NewAwsS3BucketDeprecatedServerSideEncryptionConfigurationRule(),
				rules.NewAwsS3BucketDeprecatedVersioningRule(),
				rules.NewAwsS3BucketDeprecatedWebsiteRule(),
				rules.NewAwsSecurityGroupInvalidProtocolRule(),
				rules.NewAwsSQSQueueFIFONameSuffixRule(),
			},
		},
	})
}
