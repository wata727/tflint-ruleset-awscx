# Rule Backlog

This file tracks rule ideas for this repository.

Use it to keep a durable record of:

1. candidate rules worth considering
2. source URLs
3. why a rule matters
4. implementation difficulty
5. false-positive risk
6. current decision status

Update this file when:

1. a new candidate is found
2. a candidate is selected for implementation
3. a candidate is rejected or deferred
4. a rule is implemented

## Status Guide

Use one of these statuses:

1. `new`
2. `investigating`
3. `selected`
4. `implemented`
5. `deferred`
6. `rejected`

## Candidate Template

Copy this section for each rule candidate.

```md
## <rule_name>

- Status: new
- Resource(s):
- Short description:
- Why it matters:
- Detection approach:
- False-positive risk: low | medium | high
- Implementation difficulty: low | medium | high
- Overlap notes:
- Selected on:
- Implemented on:

### Sources

- AWS docs: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_StreamSpecification.html
- Terraform Registry docs:
- terraform-provider-aws issue/PR:

### Notes

- 
```

## Candidates

## awscx_ecs_service_daemon_fargate_launch_type

- Status: implemented
- Resource(s): `aws_ecs_service`
- Short description: Disallow `scheduling_strategy = "DAEMON"` when `launch_type = "FARGATE"`.
- Why it matters: ECS does not support the `DAEMON` scheduling strategy for Fargate tasks, so this catches a concrete invalid service configuration before apply.
- Detection approach: Evaluate explicit `launch_type` and `scheduling_strategy` arguments on `aws_ecs_service` and report only when both resolve to `FARGATE` and `DAEMON`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: This complements the existing ECS daemon deployment-percent rule with another service-specific validity check and avoids the heavier overlap of Route 53 alias argument conflicts.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service_definition_parameters.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because both the provider documentation and ECS documentation describe this combination as unsupported rather than advisory guidance.
- The rule intentionally skips unknown expressions and omitted arguments to avoid guessing about capacity provider behavior.

## awscx_ecs_service_daemon_unsupported_deployment_controller

- Status: implemented
- Resource(s): `aws_ecs_service`
- Short description: Disallow `scheduling_strategy = "DAEMON"` when `deployment_controller.type` is `CODE_DEPLOY` or `EXTERNAL`.
- Why it matters: The provider documentation explicitly states that services using the `CODE_DEPLOY` or `EXTERNAL` deployment controller types do not support the `DAEMON` scheduling strategy, so this catches a concrete invalid ECS service combination before apply.
- Detection approach: Evaluate explicit `scheduling_strategy` and `deployment_controller.type` values on `aws_ecs_service` and report only when they resolve to `DAEMON` plus `CODE_DEPLOY` or `EXTERNAL`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing Fargate daemon rule and deployment percent rule by covering the remaining provider-documented unsupported daemon combinations without depending on runtime service state.
- Selected on: 2026-03-24
- Implemented on: 2026-03-24

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_CreateService.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation describes the combination as unsupported rather than optional guidance.
- The rule only reports explicit known values and skips omitted or unknown deployment controller types to avoid speculative findings.

## awscx_sqs_queue_policy_missing_current_version

- Status: implemented
- Resource(s): `aws_sqs_queue_policy`
- Short description: Require SQS queue policies to set top-level `Version = "2012-10-17"`.
- Why it matters: The AWS provider documentation warns that SQS queue create and update operations can hang indefinitely when the associated queue policy does not explicitly set the current policy language version, so this catches a concrete provider-specific failure mode before apply.
- Detection approach: Evaluate `aws_sqs_queue_policy.policy` when it resolves to a known string, parse the JSON document, and report policies whose top-level `Version` field is missing or not equal to `2012-10-17`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Adds SQS coverage for a provider-documented timeout pitfall and stays conservative by skipping unknown or invalid JSON that Terraform will already reject elsewhere.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_version.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue_policy
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation describes a concrete timeout failure mode rather than advisory guidance.
- The rule only inspects policies that evaluate to a known JSON string and skips unresolved expressions to avoid speculative findings.

## awscx_s3_bucket_deprecated_grant

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn on deprecated inline `grant` blocks on `aws_s3_bucket`.
- Why it matters: The AWS provider has deprecated inline S3 bucket grant blocks in favor of the standalone ACL resource, so keeping them in `aws_s3_bucket` creates avoidable upgrade friction without adding AWS-specific value.
- Detection approach: Inspect `aws_s3_bucket` resources and report explicit `grant` blocks.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Extends the existing S3 split-resource migration rules with another explicit inline block deprecation and shares the same conservative detection shape.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433

### Notes

- Implemented as `WARNING` because this is provider migration guidance rather than a hard API validation error.
- The rule only reports explicit inline `grant` blocks on `aws_s3_bucket` and does not inspect standalone ACL resources.

## awscx_ecs_service_deployment_maximum_percent_daemon

- Status: implemented
- Resource(s): `aws_ecs_service`
- Short description: Disallow `deployment_maximum_percent` when the ECS service uses `scheduling_strategy = "DAEMON"`.
- Why it matters: The provider documentation marks `deployment_maximum_percent` as not valid for daemon-scheduled services, so keeping it on those services creates a concrete invalid argument combination instead of expressing useful intent.
- Detection approach: Evaluate `scheduling_strategy` and report explicit `deployment_maximum_percent` arguments when it resolves to `DAEMON`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the ECS grace-period rule with another provider-documented applicability check without depending on task-definition or deployment-controller context.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service_definition_parameters.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation describes this argument as not valid for `DAEMON`, which is stronger than a recommendation.
- The rule only fires for explicit `scheduling_strategy = "DAEMON"` and skips unknown expressions to avoid guessing about module inputs.

## awscx_ecs_service_health_check_grace_period_without_load_balancer

- Status: implemented
- Resource(s): `aws_ecs_service`
- Short description: Disallow `health_check_grace_period_seconds` unless the ECS service has at least one `load_balancer` block.
- Why it matters: The provider documentation says the grace period is only valid for services configured to use load balancers, so keeping it on non-load-balanced services is a concrete invalid argument combination with a straightforward fix.
- Detection approach: Inspect `aws_ecs_service` resources and report explicit `health_check_grace_period_seconds` arguments when no `load_balancer` block is present.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Adds ECS coverage for a provider-documented validity constraint and does not overlap the existing ELB, EC2, RDS, or S3 rules.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_services.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation describes the argument as valid only with load balancers, not as advisory guidance.
- The rule checks only the explicit presence of `health_check_grace_period_seconds` and `load_balancer`, so it stays conservative around other ECS service semantics.

## awscx_sfn_state_machine_log_destination_missing_wildcard

- Status: implemented
- Resource(s): `aws_sfn_state_machine`
- Short description: Require `logging_configuration.log_destination` to end with `:*`.
- Why it matters: AWS Step Functions requires CloudWatch Logs destination ARNs to include the trailing wildcard suffix, so omitting it produces an invalid logging configuration with a clear service-specific fix.
- Detection approach: Inspect `logging_configuration.log_destination` on `aws_sfn_state_machine` and report explicit string values that do not end with `:*`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Adds Step Functions coverage for a concrete API requirement and does not overlap the existing ELB, RDS, or S3 checks.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/step-functions/latest/apireference/API_CloudWatchLogsLogGroup.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sfn_state_machine
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/39558

### Notes

- Implemented as `ERROR` because the Step Functions API documents the suffix as a hard ARN requirement rather than optional guidance.
- The rule intentionally skips unknown expressions to avoid speculative findings on variable-driven destinations.

## awscx_db_instance_blue_green_update_without_backup_retention

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Require `backup_retention_period > 0` when `blue_green_update.enabled = true`.
- Why it matters: The provider documentation requires automated backups for low-downtime updates and RDS Blue/Green deployments. Leaving backups disabled causes a concrete provider-side configuration error for an explicit upgrade path.
- Detection approach: Inspect `blue_green_update.enabled` and report when it resolves to `true` while `backup_retention_period` is omitted or resolves to `0` or less.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements existing RDS dependency rules by covering a newer Blue/Green prerequisite rather than another storage-specific constraint.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/blue-green-deployments-creating.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/38733

### Notes

- Implemented as `ERROR` because the provider documents backup retention as a prerequisite for Blue/Green usage rather than a recommendation.
- The rule only reports explicit `blue_green_update.enabled = true` and skips unknown expressions to avoid speculative findings.

## awscx_s3_bucket_configuration_expected_bucket_owner_deprecated

- Status: implemented
- Resource(s): `aws_s3_bucket_abac`, `aws_s3_bucket_accelerate_configuration`, `aws_s3_bucket_acl`, `aws_s3_bucket_cors_configuration`, `aws_s3_bucket_lifecycle_configuration`, `aws_s3_bucket_logging`, `aws_s3_bucket_metadata_configuration`, `aws_s3_bucket_object_lock_configuration`, `aws_s3_bucket_request_payment_configuration`, `aws_s3_bucket_server_side_encryption_configuration`, `aws_s3_bucket_versioning`, `aws_s3_bucket_website_configuration`
- Short description: Warn on deprecated `expected_bucket_owner` across S3 bucket configuration sub-resources.
- Why it matters: The provider has deprecated this argument across S3 bucket configuration sub-resources because the expected owner is always the current account for these APIs, so keeping it in configuration adds upgrade friction without adding protection.
- Detection approach: Inspect each affected S3 bucket configuration sub-resource and report any explicit `expected_bucket_owner` argument.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing S3 deprecation rules while covering a newer provider-wide deprecation pattern rather than another inline `aws_s3_bucket` block.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonS3/latest/userguide/logging-with-S3.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_logging
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_versioning
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/pull/46262

### Notes

- Implemented as `WARNING` because this is provider migration guidance rather than a hard API validation failure.
- The rule only reports explicit configuration on the resource types called out by the provider deprecation work.

## awscx_lb_target_group_matcher_non_http_health_check

- Status: implemented
- Resource(s): `aws_lb_target_group`
- Short description: Disallow `health_check.matcher` unless `health_check.protocol` is `HTTP` or `HTTPS`.
- Why it matters: The provider documentation limits custom matcher usage to HTTP or HTTPS health checks, with a separate allowance for Lambda target groups. Leaving `matcher` on TCP health checks is a concrete provider-level misconfiguration that commonly shows up on NLB target groups.
- Detection approach: Inspect `health_check` blocks and report `matcher` when `health_check.protocol` is explicitly set to a non-HTTP value, while skipping Lambda target groups and unknown expressions.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Adds ELB target group validation coverage without overlapping the existing listener-focused rules.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/elasticloadbalancing/latest/network/target-group-health-checks.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/8305

### Notes

- Implemented as `ERROR` because the provider documents this as an invalid attribute combination rather than a recommendation.
- The rule only reports explicit non-HTTP/HTTPS health check protocols to avoid guessing from target group defaults.

## awscx_cloudwatch_log_group_delivery_retention_in_days

- Status: implemented
- Resource(s): `aws_cloudwatch_log_group`
- Short description: Warn when `retention_in_days` is set together with `log_group_class = "DELIVERY"`.
- Why it matters: The provider documentation states that `retention_in_days` is ignored for `DELIVERY` log groups and CloudWatch Logs forces retention to 2 days, so keeping the argument in configuration creates misleading intent and can contribute to provider-side failures around delivery-class usage.
- Detection approach: Evaluate `log_group_class` and report when it resolves to `DELIVERY` while `retention_in_days` is explicitly configured.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Adds CloudWatch Logs coverage for a recent provider-specific edge case without overlapping existing EC2, RDS, or S3 rules.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/WhatIsCloudWatchLogs.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/42657

### Notes

- Implemented as a `WARNING` because the configuration is misleading and can cause drift or provider trouble, but the provider documentation describes the argument as ignored rather than as a hard validation error.

## awscx_lb_listener_alpn_policy_non_tls

- Status: implemented
- Resource(s): `aws_lb_listener`
- Short description: Disallow `alpn_policy` unless the listener protocol is `TLS`.
- Why it matters: The provider and ELB API documentation both scope ALPN policy support to TLS listeners, so attaching it to another listener protocol is an invalid configuration rather than a preference.
- Detection approach: Evaluate `protocol` and report when it resolves to a non-`TLS` value while `alpn_policy` is explicitly configured.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Another ELB listener dependency rule, but it covers a separate protocol-specific argument from the existing certificate and SSL policy checks.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_CreateListener.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider docs describe `alpn_policy` as valid only for TLS listeners.
- The rule intentionally skips listeners with unknown `protocol` expressions because the final protocol cannot be resolved statically.

## awscx_lb_listener_missing_certificate_arn

- Status: implemented
- Resource(s): `aws_lb_listener`
- Short description: Require `certificate_arn` when the listener protocol is `HTTPS` or `TLS`.
- Why it matters: The provider documentation requires a default certificate for encrypted listeners, so omitting it is almost certainly an invalid configuration.
- Detection approach: Evaluate `protocol` and report when it resolves to `HTTPS` or `TLS` but `certificate_arn` is not configured.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Focused ELB listener validity rule with minimal overlap with existing EC2 and S3 checks.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_CreateListener.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because encrypted listeners require a default certificate rather than merely recommending one.

## awscx_instance_deprecated_network_interface

- Status: implemented
- Resource(s): `aws_instance`
- Short description: Warn when deprecated `network_interface` blocks are used on `aws_instance`.
- Why it matters: The provider explicitly deprecates this block and points users to `primary_network_interface` for the primary ENI and `aws_network_interface_attachment` for additional ENIs, so continuing to use it increases upgrade friction and can lock users into replacement-prone boot-time attachment behavior.
- Detection approach: Flag any `network_interface` block present on `aws_instance`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Another provider deprecation rule, but narrowly scoped to explicit block usage with clear replacement guidance.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs:
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance
- terraform-provider-aws issue/PR:

### Notes

- Implemented as a `WARNING` because deprecated configurations can still exist during migration, but the replacement path is explicit in the provider docs.

## awscx_db_instance_storage_throughput_non_gp3

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Disallow `storage_throughput` unless `storage_type = "gp3"`.
- Why it matters: The provider docs state that `storage_throughput` can only be set with gp3 storage, so mismatched usage is almost certainly invalid.
- Detection approach: Report when `storage_throughput` is present and either `storage_type` is omitted or it evaluates to a value other than `gp3`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing `awscx_db_instance_missing_iops` rule.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#gp3-storage
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Implemented as an `ERROR` because the provider documents this as an invalid attribute combination rather than a best-practice recommendation.
- The rule intentionally skips ambiguous `storage_type` expressions and only reports explicit non-`gp3` values or omitted `storage_type`.

## awscx_db_instance_dedicated_log_volume_non_io1_io2

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Disallow `dedicated_log_volume` unless `storage_type` is `io1` or `io2`.
- Why it matters: The provider docs require dedicated log volumes to use Provisioned IOPS storage, which maps to `io1` or `io2` in Terraform.
- Detection approach: Evaluate `dedicated_log_volume` and report when it is explicitly `true` while `storage_type` is omitted or clearly not `io1`/`io2`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing RDS storage validity rules with another explicit provider-side attribute constraint.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_PIOPS.StorageTypes.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documents this as an invalid storage-type combination rather than a best-practice recommendation.
- The rule intentionally reports only explicit `dedicated_log_volume = true` values and skips unknown expressions.

## awscx_db_instance_manage_master_user_password_conflict

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Disallow `password` and `password_wo` when `manage_master_user_password = true`.
- Why it matters: RDS-managed master credentials and user-supplied master passwords are mutually exclusive, so configuring both is a concrete provider-level mistake.
- Detection approach: Evaluate `manage_master_user_password` and report any explicit `password` or `password_wo` attribute when it resolves to `true`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing RDS storage checks with a credential-management validity rule grounded in the provider docs.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-secrets-manager.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation explicitly marks these attribute combinations as invalid.
- The rule intentionally skips ambiguous expressions and only reports when `manage_master_user_password` is explicitly `true`.

## awscx_db_instance_performance_insights_arguments_without_enabled

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Disallow `performance_insights_kms_key_id` and `performance_insights_retention_period` unless `performance_insights_enabled = true`.
- Why it matters: The provider documentation makes both arguments dependent on Performance Insights being enabled, so setting them without enabling the feature is a concrete configuration mistake.
- Detection approach: Report each dependent argument when it is configured and `performance_insights_enabled` is omitted or explicitly `false`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Consolidates two closely related dependency checks into one rule to keep the ruleset smaller and the messaging more coherent.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_PerfInsights.Enabling.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation explicitly requires `performance_insights_enabled = true` when either dependent argument is set.
- The rule emits separate issues for the KMS key and retention arguments so the fix is obvious even when both appear in the same resource.

## awscx_db_instance_performance_insights_kms_key_without_enabled

- Status: rejected
- Resource(s): `aws_db_instance`
- Short description: Require `performance_insights_enabled = true` when `performance_insights_kms_key_id` is set.
- Why it matters: The provider docs require Performance Insights to be enabled before a KMS key can be set for that feature.
- Detection approach: Report when `performance_insights_kms_key_id` is configured and `performance_insights_enabled` is omitted or explicitly `false`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Strong follow-up candidate, but slightly less urgent than the selected credential-conflict rule.
- Selected on:
- Implemented on:

### Sources

- AWS docs:
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Replaced by the broader implemented rule `awscx_db_instance_performance_insights_arguments_without_enabled`.

## awscx_db_instance_performance_insights_retention_without_enabled

- Status: rejected
- Resource(s): `aws_db_instance`
- Short description: Require `performance_insights_enabled = true` when `performance_insights_retention_period` is set.
- Why it matters: The provider docs treat custom Performance Insights retention as dependent on the feature being enabled first.
- Detection approach: Report when `performance_insights_retention_period` is configured and `performance_insights_enabled` is omitted or explicitly `false`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Similar to the KMS-key candidate and likely worth grouping into a later Performance Insights-focused cycle.
- Selected on:
- Implemented on:

### Sources

- AWS docs:
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Replaced by the broader implemented rule `awscx_db_instance_performance_insights_arguments_without_enabled`.

## awscx_db_instance_database_insights_advanced_requirements

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Require `performance_insights_enabled = true` and a long enough retention period when `database_insights_mode = "advanced"`.
- Why it matters: AWS documents advanced Database Insights mode as depending on Performance Insights and a retention period of at least 465 days.
- Detection approach: Report when `database_insights_mode` resolves to `advanced` and the dependent Performance Insights settings are omitted or explicitly invalid.
- False-positive risk: low
- Implementation difficulty: medium
- Overlap notes: Adjacent to the implemented Performance Insights dependency rule, but slightly broader because it depends on newer Database Insights semantics.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_DatabaseInsights.TurningOnAdvanced.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/41607

### Notes

- Implemented as `ERROR` because AWS documents advanced mode as requiring both Performance Insights and a retention period of at least 465 days.
- The rule intentionally checks only explicit `database_insights_mode = "advanced"` values and skips ambiguous expressions for the dependent attributes.

## awscx_lb_listener_missing_ssl_policy

- Status: implemented
- Resource(s): `aws_lb_listener`
- Short description: Require `ssl_policy` when `protocol` is `HTTPS` or `TLS`.
- Why it matters: The provider documentation marks `ssl_policy` as required for encrypted listeners, making omission an explicit invalid configuration.
- Detection approach: Evaluate `protocol` and report when it resolves to `HTTPS` or `TLS` but `ssl_policy` is not configured.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Strong companion to the existing certificate ARN listener rule, but deferred to avoid clustering too many ELB listener checks in a single cycle.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_CreateListener.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because encrypted listeners without an explicit SSL policy violate a documented provider requirement rather than a best-practice recommendation.
- The rule intentionally checks only explicit `protocol` values and skips ambiguous expressions to avoid speculative reporting.

## awscx_spot_instance_request_legacy_api

- Status: deferred
- Resource(s): `aws_spot_instance_request`
- Short description: Warn on `aws_spot_instance_request` because AWS strongly discourages the legacy Spot request APIs used by this resource.
- Why it matters: AWS guidance recommends newer EC2 instance configuration paths over the legacy Spot Instance Request APIs.
- Detection approach: Flag any use of `aws_spot_instance_request`.
- False-positive risk: medium
- Implementation difficulty: low
- Overlap notes: Grounded in AWS guidance, but broader and more policy-like than explicit argument-conflict checks.
- Selected on:
- Implemented on:

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-best-practices.html#which-spot-request-method-to-use
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/spot_instance_request
- terraform-provider-aws issue/PR:

### Notes

- Deferred because warning on the entire resource would likely be noisier for teams that still intentionally use Spot Instance Requests.

## awscx_dynamodb_table_invalid_stream_view_type

- Status: implemented
- Resource(s): `aws_dynamodb_table`
- Short description: Disallow invalid `stream_view_type` values.
- Why it matters: The provider accepts only a fixed set of enum-like values, so invalid strings are almost certainly mistakes.
- Detection approach: Evaluate `stream_view_type` and compare against the allowed values.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Focused AWS/provider validation rule.
- Selected on:
- Implemented on:

### Sources

- AWS docs: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_StreamSpecification.html
- Terraform Registry docs:
- terraform-provider-aws issue/PR:

### Notes

- Add concrete source links when the exact references used for this rule are identified.

## awscx_security_group_invalid_protocol

- Status: implemented
- Resource(s): `aws_security_group`
- Short description: Disallow invalid `protocol` values in `ingress` and `egress`.
- Why it matters: Invalid protocol strings are likely configuration mistakes and should be caught before apply.
- Detection approach: Evaluate `protocol` and allow known symbolic protocol names plus numeric protocol values.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Focused provider-facing validation rule.
- Selected on:
- Implemented on:

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
- Terraform Registry docs:
- terraform-provider-aws issue/PR:

### Notes

- AWS reference aligned with the rule `Link()` implementation and `README.md`.

## awscx_s3_bucket_deprecated_acl

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `acl` is used on `aws_s3_bucket`.
- Why it matters: Inline bucket ACL usage is deprecated in the provider, and newer S3 defaults make ACL-driven access patterns more error-prone for new buckets.
- Detection approach: Flag the `acl` attribute on `aws_s3_bucket` and direct users toward the modern S3 access control resources and policies.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: This is partly a deprecation/upgrade rule rather than a pure correctness rule, but it is strongly grounded in AWS/provider behavior changes.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://aws.amazon.com/about-aws/whats-new/2022/12/amazon-s3-automatically-enable-block-public-access-disable-access-control-lists-buckets-april-2023/
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_ownership_controls
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/28353

### Notes

- Implemented as a `WARNING` to surface the provider deprecation without over-claiming that every existing bucket configuration is immediately invalid.
- The rule intentionally checks only explicit inline `acl` usage and does not try to infer bucket age or object ownership settings.

## awscx_s3_bucket_deprecated_versioning

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `versioning` is used on `aws_s3_bucket`.
- Why it matters: The AWS provider deprecated inline bucket versioning management in favor of the standalone `aws_s3_bucket_versioning` resource, so continuing to configure it inline increases upgrade friction.
- Detection approach: Flag any `versioning` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Another provider deprecation rule in the S3 bucket split-resource family, but still narrowly scoped to explicit inline usage.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- HashiCorp docs: https://developer.hashicorp.com/validated-patterns/terraform/upgrade-terraform-provider
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_versioning
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433

### Notes

- Implemented as a `WARNING` because deprecated inline configuration can still exist in older modules while users migrate to the standalone versioning resource.
- The rule intentionally reports only the presence of the inline block and does not try to validate the versioning settings themselves.

## awscx_s3_bucket_deprecated_server_side_encryption_configuration

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `server_side_encryption_configuration` is used on `aws_s3_bucket`.
- Why it matters: This is part of the same S3 bucket split-resource migration as inline versioning and would help users move to `aws_s3_bucket_server_side_encryption_configuration`.
- Detection approach: Flag any `server_side_encryption_configuration` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Sibling rule to inline versioning deprecation, still narrowly scoped to explicit inline usage.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- HashiCorp docs: https://developer.hashicorp.com/validated-patterns/terraform/upgrade-terraform-provider
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_server_side_encryption_configuration
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433

### Notes

- Implemented as a `WARNING` because deprecated inline configuration can still exist in older modules while users migrate to the standalone encryption resource.
- The rule intentionally reports only the presence of the inline block and does not attempt to validate encryption settings.

## awscx_s3_bucket_deprecated_logging

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `logging` is used on `aws_s3_bucket`.
- Why it matters: This is another S3 bucket argument included in the provider's split-resource deprecation plan.
- Detection approach: Flag any `logging` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Similar to the selected S3 deprecation rules; deferred to avoid shipping too many nearly identical warnings in one pass.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- HashiCorp docs: https://developer.hashicorp.com/validated-patterns/terraform/upgrade-terraform-provider
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_logging
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433

### Notes

- Implemented as a `WARNING` because deprecated inline configuration can still exist in older modules while users migrate to the standalone logging resource.
- The rule intentionally reports only the presence of the inline block and does not attempt to validate server access logging settings.

## awscx_s3_bucket_deprecated_lifecycle_rule

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `lifecycle_rule` is used on `aws_s3_bucket`.
- Why it matters: The AWS provider split S3 bucket lifecycle management into `aws_s3_bucket_lifecycle_configuration`, so keeping lifecycle rules inline increases upgrade friction and drifts away from current provider guidance.
- Detection approach: Flag any `lifecycle_rule` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Another S3 split-resource deprecation rule, but still narrowly scoped to explicit inline usage with a single replacement resource.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs:
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433
- HashiCorp blog: https://www.hashicorp.com/blog/terraform-aws-provider-4-0-refactors-s3-bucket-resource

### Notes

- Implemented as a `WARNING` because deprecated inline configuration can still exist in older modules while users migrate to the standalone lifecycle configuration resource.
- The rule intentionally reports only the presence of the inline block and does not try to validate individual lifecycle sub-block semantics.

## awscx_s3_bucket_deprecated_replication_configuration

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `replication_configuration` is used on `aws_s3_bucket`.
- Why it matters: The provider moved S3 bucket replication management to `aws_s3_bucket_replication_configuration`, so inline configuration increases upgrade friction.
- Detection approach: Flag any `replication_configuration` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Another S3 split-resource deprecation rule, but still narrowly scoped to explicit inline usage with a single replacement resource.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs:
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_replication_configuration
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433
- HashiCorp blog: https://www.hashicorp.com/blog/terraform-aws-provider-4-0-refactors-s3-bucket-resource

### Notes

- Implemented as a `WARNING` because deprecated inline configuration can still exist in older modules while users migrate to the standalone replication configuration resource.
- The rule intentionally reports only the presence of the inline block and does not attempt to validate replication sub-block semantics.

## awscx_s3_bucket_deprecated_website

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `website` is used on `aws_s3_bucket`.
- Why it matters: The provider moved static website configuration management to `aws_s3_bucket_website_configuration`, so keeping website settings inline adds upgrade friction and diverges from the current resource model.
- Detection approach: Flag any `website` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Another S3 split-resource deprecation rule, but still narrowly scoped to explicit inline usage with a single replacement resource.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteHosting.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433

### Notes

- Implemented as a `WARNING` because deprecated inline website configuration can still exist in older modules while users migrate to the standalone website configuration resource.
- The rule intentionally reports only the presence of the inline block and does not attempt to validate website sub-block semantics.

## awscx_db_instance_publicly_accessible

- Status: deferred
- Resource(s): `aws_db_instance`
- Short description: Warn when `publicly_accessible = true` is set on an RDS instance.
- Why it matters: Publicly reachable databases increase exposure risk when not strongly justified.
- Detection approach: Flag an explicit `publicly_accessible = true`.
- False-positive risk: medium
- Implementation difficulty: low
- Overlap notes: Security value is real, but deployment intent varies enough that it is noisier than the selected deprecation candidate.
- Selected on:
- Implemented on:

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_VPC.WorkingWithRDSInstanceinaVPC.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Deferred because the rule would encode a policy preference rather than a provider or AWS-side validity requirement.

## awscx_lb_deletion_protection_disabled

- Status: deferred
- Resource(s): `aws_lb`
- Short description: Warn when `enable_deletion_protection = false` is explicitly set.
- Why it matters: Deletion protection can reduce the blast radius of accidental load balancer removal.
- Detection approach: Flag an explicit `enable_deletion_protection = false`.
- False-positive risk: medium
- Implementation difficulty: low
- Overlap notes: Operationally useful, but still closer to organization policy than a provider-side requirement.
- Selected on:
- Implemented on:

### Sources

- AWS docs: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/edit-load-balancer-attributes.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb
- terraform-provider-aws issue/PR:

### Notes

- Deferred because ephemeral or lower-environment load balancers may intentionally disable deletion protection.

## awscx_eks_addon_deprecated_resolve_conflicts

- Status: implemented
- Resource(s): `aws_eks_addon`
- Short description: Warn when deprecated `resolve_conflicts` is used.
- Why it matters: The provider deprecated `resolve_conflicts` in favor of separate create and update attributes because create-time and update-time semantics differ, especially for `PRESERVE`.
- Detection approach: Flag any explicit `resolve_conflicts` attribute on `aws_eks_addon`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Provider deprecation rule with direct user migration guidance; intentionally limited to explicit attribute usage.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateAddon.html
- AWS docs: https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateAddon.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_addon
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/27481

### Notes

- Implemented as a `WARNING` because the configuration may still work for some values, but it is deprecated and the replacement attributes are more precise.
- The rule does not inspect the attribute value because any explicit use of `resolve_conflicts` should migrate to the split create/update attributes.

## awscx_launch_template_imdsv2_optional_tokens

- Status: implemented
- Resource(s): `aws_launch_template`
- Short description: Warn when `metadata_options.http_tokens = "optional"` explicitly allows IMDSv1.
- Why it matters: AWS recommends requiring IMDSv2 for new instances, and launch-template metadata settings are a common place to enforce that defense-in-depth control.
- Detection approach: Inspect `metadata_options.http_tokens` on `aws_launch_template` and report only the explicit `"optional"` setting.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: A narrower version of the original IMDSv2 candidate, chosen to avoid guessing account-level defaults, AMI settings, or organization policies.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/25909

### Notes

- Implemented as a `WARNING` because some environments may intentionally retain IMDSv1 compatibility for older software.
- The rule does not flag omitted `metadata_options`, because account-level IMDS defaults and AMI configuration can make omission safe.

## awscx_instance_imdsv2_optional_tokens

- Status: implemented
- Resource(s): `aws_instance`
- Short description: Warn when `metadata_options.http_tokens = "optional"` explicitly allows IMDSv1.
- Why it matters: AWS recommends requiring IMDSv2 for new instances, and `aws_instance` is a common direct entry point for EC2 configuration.
- Detection approach: Inspect `metadata_options.http_tokens` on `aws_instance` and report only the explicit `"optional"` setting.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Parallel coverage for direct EC2 instances; kept narrow for the same reasons as the launch-template rule.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/10949

### Notes

- Implemented as a `WARNING` because some environments may still depend on IMDSv1 compatibility during migration.
- The rule does not flag omitted `metadata_options`, because account-level defaults can already require IMDSv2.

## awscx_efs_file_system_missing_provisioned_throughput

- Status: implemented
- Resource(s): `aws_efs_file_system`
- Short description: Require `provisioned_throughput_in_mibps` when `throughput_mode = "provisioned"`.
- Why it matters: AWS requires a provisioned throughput value when EFS is configured for provisioned throughput mode.
- Detection approach: Evaluate `throughput_mode` and report when it is explicitly `provisioned` while `provisioned_throughput_in_mibps` is absent.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Focused provider-facing validity check with a direct AWS-side requirement.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/efs/latest/ug/throughput-modes.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the configuration is incomplete when provisioned mode is selected without a throughput value.
- The rule only checks the missing-attribute case and does not validate numeric limits, which may vary by Region and quota.

## awscx_efs_file_system_provisioned_throughput_non_provisioned

- Status: implemented
- Resource(s): `aws_efs_file_system`
- Short description: Disallow `provisioned_throughput_in_mibps` unless `throughput_mode = "provisioned"`.
- Why it matters: The provider docs mark provisioned throughput as applicable only with `throughput_mode = "provisioned"`, so keeping the attribute on `bursting`, `elastic`, or omitted mode settings creates an invalid or misleading configuration.
- Detection approach: Report when `provisioned_throughput_in_mibps` is present and `throughput_mode` is omitted or evaluates to a value other than `provisioned`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Direct companion to `awscx_efs_file_system_missing_provisioned_throughput`, covering the inverse argument constraint without widening scope beyond an explicit provider rule.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/efs/latest/ug/throughput-modes.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider describes `provisioned_throughput_in_mibps` as valid only with provisioned throughput mode.
- The rule intentionally skips unknown `throughput_mode` expressions and only reports omitted or explicit non-`provisioned` values.

## awscx_efs_file_system_kms_key_without_encrypted

- Status: implemented
- Resource(s): `aws_efs_file_system`
- Short description: Disallow `kms_key_id` unless `encrypted = true`.
- Why it matters: The provider docs state that custom KMS keys apply only to encrypted EFS file systems, so setting `kms_key_id` without encryption is a direct configuration error rather than a policy preference.
- Detection approach: Report when `kms_key_id` is present and `encrypted` is omitted or evaluates to `false`, while skipping unknown expressions.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing EFS throughput dependency rules with another explicit provider-side prerequisite on the same resource.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/efs/latest/ug/encryption-at-rest.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because `kms_key_id` without encryption conflicts with the documented resource contract.
- The rule intentionally checks only explicit `kms_key_id` usage and does not validate the referenced key itself.

## awscx_s3_bucket_deprecated_acceleration_status

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn on deprecated inline `acceleration_status` on `aws_s3_bucket`.
- Why it matters: The provider deprecates inline acceleration configuration in favor of `aws_s3_bucket_accelerate_configuration`, so the rule would be low-noise and easy to explain.
- Detection approach: Flag any explicit `acceleration_status` attribute on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Fits the existing S3 split-resource migration pattern without duplicating the already implemented inline `acl`, `logging`, `versioning`, or `website` checks.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonS3/latest/userguide/transfer-acceleration.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/s3_bucket.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as a `WARNING` because this is provider deprecation and migration guidance rather than a hard API validation failure.
- The rule intentionally reports only explicit inline `acceleration_status` usage and points users at the standalone accelerate configuration resource.

## awscx_s3_bucket_deprecated_policy

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn on deprecated inline `policy` on `aws_s3_bucket`.
- Why it matters: The provider deprecates inline bucket policy management in favor of `aws_s3_bucket_policy`, and the inline form can also produce confusing perpetual diff behavior when the policy is not specific enough.
- Detection approach: Flag any explicit `policy` attribute on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Fits the existing S3 split-resource migration pattern without overlapping the already implemented inline deprecation rules for ACL, logging, lifecycle, replication, encryption, versioning, website, and acceleration.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonS3/latest/userguide/example-bucket-policies.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/s3_bucket.html.markdown
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_policy

### Notes

- Implemented as a `WARNING` because the provider docs frame this as deprecation and migration guidance rather than an API-level invalid configuration.
- The rule only reports explicit inline `policy` usage and directs users to `aws_s3_bucket_policy`.

## awscx_ebs_volume_missing_iops

- Status: implemented
- Resource(s): `aws_ebs_volume`
- Short description: Require `iops` when `type = "io1"` or `type = "io2"`.
- Why it matters: AWS requires an IOPS value for provisioned IOPS EBS volume types.
- Detection approach: Evaluate `type` and report when it is explicitly `io1` or `io2` while `iops` is absent.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Focused validity check for a direct EBS API requirement.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreateVolume.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the configuration is incomplete for provisioned IOPS volume types without an explicit IOPS value.
- The rule does not validate the allowed numeric range for `iops`, which depends on volume type and instance support.

## awscx_ebs_volume_throughput_non_gp3

- Status: implemented
- Resource(s): `aws_ebs_volume`
- Short description: Disallow `throughput` unless `type = "gp3"`.
- Why it matters: The provider and EC2 API both document throughput as valid only for gp3 EBS volumes, so using it elsewhere is almost certainly an invalid configuration.
- Detection approach: Report when `throughput` is present and either `type` is omitted or it evaluates to a value other than `gp3`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Sibling validity rule to `awscx_ebs_volume_missing_iops`, but focused on the gp3-only throughput setting.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_EbsBlockDevice.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because `throughput` outside `gp3` is documented as invalid rather than advisory guidance.
- The rule intentionally skips unknown `type` expressions and only reports explicit non-`gp3` values or omitted `type`.

## awscx_db_instance_missing_iops

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Require `iops` when `storage_type = "io1"`, `storage_type = "io2"`, or `storage_type = "gp3"`.
- Why it matters: RDS requires an explicit IOPS value for these storage types.
- Detection approach: Evaluate `storage_type` and report when it is explicitly `io1`, `io2`, or `gp3` while `iops` is absent.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Direct RDS API validity check rather than an opinionated best-practice rule.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/goto/WebAPI/rds-2014-10-31/CreateDBInstance
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the configuration is incomplete when one of these storage types is selected without IOPS.
- The rule intentionally does not validate engine-specific or size-specific IOPS ranges.

## awscx_launch_template_deprecated_elastic_gpu_specifications

- Status: implemented
- Resource(s): `aws_launch_template`
- Short description: Warn when deprecated `elastic_gpu_specifications` is used.
- Why it matters: Amazon Elastic Graphics reached end of life, and the provider tracks this block for deprecation/removal.
- Detection approach: Flag any explicit `elastic_gpu_specifications` block on `aws_launch_template`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Another explicit deprecation rule, but grounded in AWS service retirement rather than a generic style preference.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_ElasticGpuSpecificationResponse.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/37589

### Notes

- Implemented as a `WARNING` because the configuration is deprecated and obsolete, but older modules may still carry the block during migration.
- The rule intentionally checks only explicit launch template usage and does not try to infer replacement behavior for other EC2 resources.

## awscx_api_gateway_deployment_deprecated_stage_management

- Status: implemented
- Resource(s): `aws_api_gateway_deployment`
- Short description: Warn when deprecated stage-management fields are used on `aws_api_gateway_deployment`.
- Why it matters: The AWS provider deprecated `stage_name`, `stage_description`, and `canary_settings` because they implicitly create or modify stages from the deployment resource, which is confusing and will be removed in a future provider version.
- Detection approach: Report explicit use of `stage_name`, `stage_description`, or `canary_settings` on `aws_api_gateway_deployment`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Narrow provider-deprecation coverage for API Gateway; it avoids broader API Gateway policy checks that would need more context.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- Terraform provider issue: https://github.com/hashicorp/terraform-provider-aws/issues/39957
- Terraform provider issue: https://github.com/hashicorp/terraform-provider-aws/issues/39958
- AWS docs: https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-deployments.html
- AWS docs: https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-stages.html

### Notes

- Implemented as a `WARNING` because the configuration still works in older provider versions, but it is explicitly deprecated and scheduled for removal.
- The rule intentionally checks only explicit deprecated deployment-side stage management and does not try to infer whether a separate `aws_api_gateway_stage` resource should exist elsewhere in the module.

## awscx_sqs_queue_fifo_name_suffix

- Status: implemented
- Resource(s): `aws_sqs_queue`
- Short description: Require FIFO queue names to end with `.fifo`.
- Why it matters: AWS requires FIFO queue names to use the `.fifo` suffix, so omitting it produces an invalid queue definition.
- Detection approach: Evaluate `fifo_queue` and report when it is explicitly `true` while an explicit `name` does not end with `.fifo`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Direct SQS API validity check; it intentionally skips `name_prefix` and omitted-name cases to avoid guessing the final generated queue name.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because FIFO queue creation requires the `.fifo` suffix.
- The rule only checks explicit `name` values and does not attempt to infer suffixes from `name_prefix` or unknown expressions.

## awscx_eip_instance_network_interface_conflict

- Status: implemented
- Resource(s): `aws_eip`
- Short description: Disallow setting both `instance` and `network_interface` on `aws_eip`.
- Why it matters: The provider documentation and EC2 AssociateAddress API both state that an Elastic IP association must target either an instance or a network interface, but not both.
- Detection approach: Report when an `aws_eip` resource explicitly sets both `instance` and `network_interface`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Direct API-validity check with no cross-resource inference or environment-specific policy assumptions.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/eip.html.markdown
- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_AssociateAddress.html

### Notes

- Implemented as `ERROR` because the combination is documented as invalid and can lead to undefined behavior rather than a style preference.
- The rule only checks explicit co-presence of the two arguments and does not attempt to validate the referenced resources.

## awscx_security_group_all_protocol_nonzero_ports

- Status: implemented
- Resource(s): `aws_security_group`
- Short description: Require inline `ingress` and `egress` rules with `protocol = "-1"` to use `from_port = 0` and `to_port = 0`.
- Why it matters: AWS treats `IpProtocol = -1` as all protocols and all ports, so keeping nonzero ports in Terraform is misleading and can mask an unintended any-port rule.
- Detection approach: Report inline security group rules where `protocol` evaluates to `"-1"` and either `from_port` or `to_port` evaluates to a nonzero number.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements `awscx_security_group_invalid_protocol` by validating a semantically dangerous all-protocol configuration rather than protocol spelling.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/security_group_rule.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider and EC2 API semantics are explicit, and a nonzero port range with `protocol = "-1"` is misleading rather than intentional precision.
- The rule intentionally skips unknown port expressions instead of speculating about variable values.

## awscx_lb_listener_https_missing_certificate_arn

- Status: deferred
- Resource(s): `aws_lb_listener`
- Short description: Report HTTPS listeners that omit `certificate_arn`.
- Why it matters: HTTPS listeners usually need an attached certificate to terminate TLS correctly.
- Detection approach: Check `protocol = "HTTPS"` and flag missing `certificate_arn`.
- False-positive risk: medium
- Implementation difficulty: low
- Overlap notes: Potentially useful, but the provider documentation needed more careful validation around protocol variants and defaults before turning it into a correctness rule.
- Selected on:
- Implemented on:

### Sources

- AWS docs:
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/lb_listener.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Deferred because the certificate requirements differ across listener protocols and load balancer types, so a narrow first implementation needs more source review.

## awscx_sqs_queue_high_throughput_fifo_partial_settings

- Status: deferred
- Resource(s): `aws_sqs_queue`
- Short description: Warn when only one of the high-throughput FIFO tuning attributes is set.
- Why it matters: Partial configuration can look like high-throughput FIFO is enabled when queue-level deduplication or throughput defaults still apply.
- Detection approach: Compare explicit `deduplication_scope` and `fifo_throughput_limit` values and warn on one-sided configuration.
- False-positive risk: medium
- Implementation difficulty: low
- Overlap notes: Useful advisory candidate, but it is closer to performance tuning guidance than a hard validity check.
- Selected on:
- Implemented on:

### Sources

- AWS docs: https://aws.amazon.com/about-aws/whats-new/2021/05/amazon-sqs-now-supports-a-high-throughput-mode-for-fifo-queues/
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/sqs_queue.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Deferred because the rule would be advisory rather than invalidity-driven, and this cycle favored a lower-noise correctness rule.

## awscx_sns_topic_fifo_name_suffix

- Status: implemented
- Resource(s): `aws_sns_topic`
- Short description: Require FIFO SNS topic names to end with `.fifo`.
- Why it matters: AWS and the provider both document the `.fifo` suffix as a hard requirement for FIFO topics, so missing it is a direct configuration error rather than a style preference.
- Detection approach: Report resources that explicitly set `fifo_topic = true` and `name` without the `.fifo` suffix.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Sibling to the existing SQS FIFO suffix rule, but still worthwhile because SNS FIFO topics use a different resource and API.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/sns/latest/api/API_CreateTopic.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/sns_topic.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as an `ERROR` because FIFO topic creation requires the `.fifo` suffix.
- The rule intentionally skips cases that use `name_prefix` or unknown expressions, because the final generated name is not statically known.

## awscx_lb_listener_mutual_authentication_verify_requirements

- Status: implemented
- Resource(s): `aws_lb_listener`
- Short description: Require `mutual_authentication.trust_store_arn` for `mode = "verify"` and disallow verify-only attributes for `off` or `passthrough`.
- Why it matters: The provider documentation explicitly marks `trust_store_arn` as required for verify mode and `advertise_trust_store_ca_names` / `ignore_client_certificate_expiry` as invalid outside verify mode. A provider issue also shows that leaving verify-only attributes in non-verify modes can produce ELB API validation failures.
- Detection approach: Evaluate `mutual_authentication.mode` and report missing `trust_store_arn` for explicit `verify` values, or any explicit verify-only attributes for explicit `off` / `passthrough` values.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Extends existing `aws_lb_listener` protocol and certificate checks with a nested-block validity rule that is still based on explicit provider contracts rather than organization-specific policy.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/en_us/elasticloadbalancing/latest/APIReference/API_MutualAuthenticationAttributes.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/lb_listener.html.markdown
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/34861

### Notes

- Implemented as `ERROR` because the documented combinations are invalid and can fail at provider or ELB API validation time.
- The rule intentionally skips unknown `mode` expressions to avoid speculative reporting on variable-driven configurations.

## awscx_autoscaling_group_invalid_max_instance_lifetime

- Status: implemented
- Resource(s): `aws_autoscaling_group`
- Short description: Disallow `max_instance_lifetime` values other than `0` or `86400..31536000`.
- Why it matters: The provider documentation defines a strict numeric range, so invalid values are rejected and can be caught statically.
- Detection approach: Evaluate explicit numeric `max_instance_lifetime` values and report anything outside the documented range except `0`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Adds Auto Scaling coverage without overlapping the existing EC2 launch template and instance rules, and stays low-noise by only validating explicit numeric values.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-max-instance-lifetime.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/autoscaling_group
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/autoscaling_group.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation defines this as a concrete validity range rather than a best-practice recommendation.
- The rule skips unknown and non-integer expressions to avoid speculative reporting on variable-driven or computed configurations.

## awscx_lb_target_group_protocol_version_non_http

- Status: implemented
- Resource(s): `aws_lb_target_group`
- Short description: Disallow `protocol_version` unless `protocol` is `HTTP` or `HTTPS`, and flag it on explicit Lambda target groups.
- Why it matters: The provider documentation marks `protocol_version` as only applicable for HTTP/HTTPS target groups, while AWS documents protocol versions as an Application Load Balancer backend feature for HTTP-family traffic. Setting it on TCP/TLS/UDP-style target groups or Lambda target groups is concrete provider misuse.
- Detection approach: Report explicit `protocol_version` attributes when `protocol` is explicitly not `HTTP` or `HTTPS`, or when `target_type` is explicitly `lambda`. Skip unknown expressions to avoid speculative issues.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing `aws_lb_target_group` health check rule with another top-level argument compatibility check based on explicit provider constraints.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-target-groups.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/lb_target_group.html.markdown
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/15929

### Notes

- Implemented as `ERROR` because the provider documentation treats this as an argument applicability constraint rather than a stylistic recommendation.
- The rule intentionally skips unknown `protocol` or `target_type` expressions unless the invalid context can be proven statically.

## awscx_db_instance_enhanced_monitoring_role_requirements

- Status: implemented
- Resource(s): `aws_db_instance`
- Short description: Require `monitoring_role_arn` when `monitoring_interval` enables Enhanced Monitoring, and disallow the role when interval is `0`.
- Why it matters: The provider documentation models `monitoring_role_arn` as dependent on Enhanced Monitoring, so mismatched settings are concrete provider misuse rather than a preference.
- Detection approach: Evaluate explicit `monitoring_interval` values and report missing or extraneous `monitoring_role_arn` based on whether monitoring is enabled.
- False-positive risk: low
- Implementation difficulty: medium
- Overlap notes: Extends the existing RDS dependency checks with another provider-documented argument pairing and stays low-noise by only acting on explicit, resolvable monitoring intervals.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.OS.Enabling.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/db_instance.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider and AWS documentation treat the role and interval as a required pair when Enhanced Monitoring is enabled.
- The rule intentionally skips configurations that set `monitoring_role_arn` without an explicit `monitoring_interval` to avoid relying on provider defaults for reporting.

## awscx_launch_template_deprecated_elastic_inference_accelerator

- Status: implemented
- Resource(s): `aws_launch_template`
- Short description: Warn on deprecated `elastic_inference_accelerator` blocks on `aws_launch_template`.
- Why it matters: Amazon Elastic Inference reached end of life in April 2024 and is no longer available, and the Terraform AWS Provider has already marked the launch template attribute for removal in the v6 line. Keeping the block in configuration creates avoidable upgrade and apply failures.
- Detection approach: Flag any explicit `elastic_inference_accelerator` block present on `aws_launch_template`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Closely related to the existing Elastic Graphics deprecation rule, but covers a separate AWS service retirement and a distinct launch template block.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_RequestLaunchTemplateData.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/41101

### Notes

- Implemented as `WARNING` because this is deprecation and service-retirement guidance rather than a new static provider validation error.
- The rule only reports explicit launch template blocks and does not try to infer behavior from other EC2 or ECS resources.

## awscx_sns_topic_fifo_attributes_without_fifo_topic

- Status: implemented
- Resource(s): `aws_sns_topic`
- Short description: Disallow FIFO-only SNS topic attributes unless `fifo_topic = true`.
- Why it matters: The provider documentation scopes `archive_policy`, `content_based_deduplication`, and `fifo_throughput_scope` to FIFO topics, so setting them on a standard topic is a concrete configuration mistake or misleading no-op.
- Detection approach: Report each explicit FIFO-only attribute when `fifo_topic` is omitted or explicitly `false`, and skip unknown expressions to avoid speculative issues.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing SNS FIFO name suffix rule with top-level attribute applicability checks on the same resource type.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/sns/latest/dg/fifo-message-dedup.html
- AWS docs: https://docs.aws.amazon.com/sns/latest/dg/message-archiving-and-replay-topic-owner.html
- AWS docs: https://docs.aws.amazon.com/sns/latest/dg/fifo-high-throughput.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/sns_topic.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the flagged attributes are documented as FIFO-topic settings rather than advisory best practices.
- The rule intentionally reports each attribute separately so the invalid configuration is obvious when multiple FIFO-only settings appear together.

## awscx_s3_bucket_deprecated_object_lock_configuration

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn on deprecated inline `object_lock_configuration` on `aws_s3_bucket`.
- Why it matters: The provider documentation deprecates inline object lock configuration in favor of `object_lock_enabled` plus `aws_s3_bucket_object_lock_configuration`, so keeping the old block adds upgrade friction and obscures the supported resource split.
- Detection approach: Flag any explicit `object_lock_configuration` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Extends the existing S3 split-resource migration rules with another explicit deprecated block, without requiring any cross-resource inference.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonS3/latest/dev/object-lock.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/s3_bucket.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `WARNING` because this is provider migration guidance rather than a new hard API validation error.
- The rule only reports the deprecated inline block and does not inspect retention settings inside the standalone replacement resource.

## awscx_s3_bucket_deprecated_request_payer

- Status: implemented
- Resource(s): `aws_s3_bucket`
- Short description: Warn on deprecated inline `request_payer` on `aws_s3_bucket`.
- Why it matters: The provider documentation deprecates inline `request_payer` in favor of `aws_s3_bucket_request_payment_configuration`, so leaving it on the bucket resource adds migration churn without adding distinct behavior.
- Detection approach: Flag any explicit `request_payer` attribute present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Continues the existing S3 split-resource migration warnings with another one-attribute deprecation check.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonS3/latest/userguide/RequesterPaysBuckets.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/s3_bucket.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `WARNING` because this is provider migration guidance rather than a hard API validation error.
- The rule only reports explicit inline usage and directs users to `aws_s3_bucket_request_payment_configuration`.

## awscx_lb_target_group_lambda_top_level_attributes

- Status: implemented
- Resource(s): `aws_lb_target_group`
- Short description: Disallow `port`, `protocol`, and `vpc_id` when `target_type = "lambda"`.
- Why it matters: The provider documentation and ELB CreateTargetGroup API both state that these top-level target-group transport settings do not apply to Lambda targets, so keeping them in configuration is a concrete invalid-context mistake.
- Detection approach: Evaluate explicit `target_type` values and report each configured top-level attribute among `port`, `protocol`, and `vpc_id` when the target type resolves to `lambda`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing `awscx_lb_target_group_protocol_version_non_http` rule by covering the other top-level attributes that become inapplicable for Lambda target groups.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/cli/latest/reference/elbv2/create-target-group.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/lb_target_group.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the flagged attributes are documented as not applying to Lambda target groups rather than being advisory best practices.
- The rule only reports explicit `target_type = "lambda"` values and skips unknown expressions to avoid speculative findings.

## awscx_cloudfront_distribution_minimum_protocol_version_default_certificate

- Status: implemented
- Resource(s): `aws_cloudfront_distribution`
- Short description: Disallow `viewer_certificate.minimum_protocol_version` when `viewer_certificate.cloudfront_default_certificate = true`.
- Why it matters: The provider documentation says `minimum_protocol_version` can only be set when `cloudfront_default_certificate = false`, so keeping it alongside the default CloudFront certificate is a concrete invalid combination rather than a stylistic preference.
- Detection approach: Inspect explicit `viewer_certificate` blocks and report `minimum_protocol_version` when `cloudfront_default_certificate` evaluates to `true`; skip omitted or unknown values to avoid speculative findings.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Complements the existing load balancer protocol applicability rules with another AWS CDN certificate applicability check grounded in provider-documented block semantics.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- AWS docs: https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/secure-connections-supported-viewer-protocols-ciphers.html
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudfront_distribution
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/cloudfront_distribution.html.markdown
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the provider documentation describes this as a direct validity constraint on the `viewer_certificate` block.
- The rule intentionally checks only the explicit `cloudfront_default_certificate = true` case and does not infer behavior when the boolean comes from unknown expressions.

## awscx_lambda_function_zip_required_attributes

- Status: implemented
- Resource(s): `aws_lambda_function`
- Short description: Require `handler` and `runtime` when `package_type` is `Zip` or omitted.
- Why it matters: The provider documentation says `handler` and `runtime` are required for Zip-based Lambda functions, while image-based functions do not need them. Omitting these fields from a Zip function produces a concrete invalid configuration rather than a stylistic disagreement.
- Detection approach: Treat omitted `package_type` as the provider default `Zip`, report missing `handler` and `runtime`, and skip unknown `package_type` expressions to avoid speculative findings.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Covers a provider-documented Lambda applicability requirement that is not already enforced by the existing ruleset and does not depend on runtime account context.
- Selected on: 2026-03-23
- Implemented on: 2026-03-23

### Sources

- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lambda_function
- Raw provider docs: https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/lambda_function.html.markdown
- AWS docs: https://docs.aws.amazon.com/lambda/latest/api/API_CreateFunction.html
- terraform-provider-aws issue/PR:

### Notes

- Implemented as `ERROR` because the missing attributes violate the documented requirements for Zip package functions.
- The rule intentionally skips unknown `package_type` expressions and only relies on the explicit `Image` exception or the documented default `Zip`.

## Backlog Hygiene

Prefer keeping this file concise and current.

When a rule is implemented:

1. keep the entry
2. mark it as `implemented`
3. fill in the sources used
4. add any follow-up ideas that came from the implementation

When a rule is rejected:

1. keep the entry if the rejection teaches something useful
2. explain the rejection briefly
3. record whether the issue was false-positive risk, weak value, or implementation complexity
