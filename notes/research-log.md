# Research Log

This file records each research and implementation cycle.

Each entry should be short, but should leave enough context for the next cycle to continue without redoing the same exploration.

## Entry Template

```md
## YYYY-MM-DD - Cycle <n>

- Goal:
- Candidates investigated:
- Selected candidate:
- Why selected:
- Sources used:
  - 
- Files changed:
  - 
- Tests run:
  - 
- Result:
- Follow-up ideas:
  - 
```

## 2026-03-22 - Repository setup

- Goal: Establish the initial repository guidance and prepare for sustained rule discovery and implementation.
- Candidates investigated:
  - `awscx_dynamodb_table_invalid_stream_view_type`
  - `awscx_security_group_invalid_protocol`
- Selected candidate:
  - Existing implemented rules were reviewed rather than newly selected in this cycle.
- Why selected:
  - This cycle focused on repository setup and validating that the current state supports continued work.
- Sources used:
  - Repository-local sources only
- Files changed:
  - `AGENTS.md`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - The repository is in a workable state for continued rule development, with room to improve durable source tracking and verification workflow.
- Follow-up ideas:
  - Add concrete source URLs to implemented rules and README entries.
  - Add a minimal TFLint verification workflow or fixture.
  - Start recording candidate comparisons in the backlog before each new implementation.

## 2026-03-22 - Operations hardening

- Goal: Make continued rule work less fragile by aligning repository guidance, durable notes, and source tracking.
- Candidates investigated:
  - No new rule candidate selected in this cycle
- Selected candidate:
  - None
- Why selected:
  - This cycle focused on operational readiness rather than new implementation work.
- Sources used:
  - https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_StreamSpecification.html
  - https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
- Files changed:
  - `AGENTS.md`
  - `README.md`
  - `rules/aws_dynamodb_table_invalid_stream_view_type.go`
  - `rules/aws_security_group_invalid_protocol.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Repository guidance now points at `go test ./...` as the minimum verification step, avoids forcing direct pushes to `main`, and keeps source URLs with the implemented rules.
- Follow-up ideas:
  - Add a minimal `.tflint.hcl` and fixture for plugin-level verification.
  - Add Terraform Registry or provider issue links when a stable, directly relevant page is identified.

## 2026-03-22 - Local verification scaffolding

- Goal: Make plugin-level verification repeatable inside the repository.
- Candidates investigated:
  - No new rule candidate selected in this cycle
- Selected candidate:
  - None
- Why selected:
  - This cycle focused on improving the verification workflow for future rule work.
- Sources used:
  - Repository-local sources only
- Files changed:
  - `.gitignore`
  - `.tflint.hcl`
  - `Makefile`
  - `README.md`
  - `AGENTS.md`
  - `testdata/tflint/invalid_security_group_protocol/main.tf`
  - `testdata/tflint/invalid_dynamodb_stream_view_type/main.tf`
- Tests run:
  - `go test ./...`
  - `make build`
  - `make verify-plugin`
- Result:
  - The repository now includes a local TFLint config, repeatable plugin verification fixtures, and a `make verify-plugin` target.
  - `make verify-plugin` completes successfully and reports the expected issues for the intentionally invalid fixtures.
- Follow-up ideas:
  - Run `make verify-plugin` after installing the plugin locally.
  - Add a valid fixture directory if a future cycle needs a no-issue plugin verification case.
  - Note: this verification scaffolding was later removed from the repository and should be treated as historical context only.

## 2026-03-22 - Backlog seeding and fixture cleanup

- Goal: Remove the last ad-hoc root fixture and seed the backlog with near-term candidates for the next work cycles.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_acl`
  - `awscx_launch_template_require_imdsv2`
- Selected candidate:
  - None
- Why selected:
  - This cycle focused on preparing cleaner inputs for the next implementation cycle instead of starting a new rule immediately.
- Sources used:
  - https://aws.amazon.com/about-aws/whats-new/2022/12/amazon-s3-automatically-enable-block-public-access-disable-access-control-lists-buckets-april-2023/
  - https://github.com/hashicorp/terraform-provider-aws/issues/28353
  - https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html
  - https://github.com/hashicorp/terraform-provider-aws/issues/25909
- Files changed:
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
  - `main.tf`
- Tests run:
  - not run
- Result:
  - The repository now keeps verification fixtures under `testdata/` only, and the backlog has two concrete candidates ready for the next implementation cycle.
- Follow-up ideas:
  - Choose whether this repository should include advisory `WARNING` rules such as deprecation and hardening guidance.
  - Add Terraform Registry documentation links for the two backlog candidates when selecting one for implementation.
  - Note: the referenced fixture cleanup applied to a previous repository state and is preserved here as history.

## 2026-03-23 - Cycle 1

- Goal: Research a practical AWS-specific rule candidate and implement one low-noise rule.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_acl`
  - `awscx_launch_template_require_imdsv2`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_s3_bucket_deprecated_acl`
- Why selected:
  - The rule is strongly grounded in provider deprecation guidance and AWS's post-April-2023 S3 defaults.
  - Detection is simple and low-noise because it only flags explicit inline `acl` usage on `aws_s3_bucket`.
  - The alternative candidates were more sensitive to environment-specific intent and would have needed looser heuristics to avoid false positives.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_ownership_controls
  - https://github.com/hashicorp/terraform-provider-aws/issues/28353
  - https://aws.amazon.com/about-aws/whats-new/2022/12/amazon-s3-automatically-enable-block-public-access-disable-access-control-lists-buckets-april-2023/
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_s3_bucket_deprecated_acl.go`
  - `rules/aws_s3_bucket_deprecated_acl_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated inline `acl` usage on `aws_s3_bucket` and points users toward the newer S3 ACL and ownership-control resources.
- Follow-up ideas:
  - Revisit launch-template IMDS checks with an explicit `http_tokens = "optional"` trigger to keep noise low.
  - Explore another S3 rule around ownership controls only if it can be implemented without cross-resource guesswork.

## 2026-03-23 - Cycle 2

- Goal: Implement another AWS-specific rule with strong security value while keeping false positives low.
- Candidates investigated:
  - `awscx_launch_template_imdsv2_optional_tokens`
  - `awscx_db_instance_publicly_accessible`
  - `awscx_instance_imdsv2_optional_tokens`
- Selected candidate:
  - `awscx_launch_template_imdsv2_optional_tokens`
- Why selected:
  - AWS directly recommends requiring IMDSv2 for new instances, and `http_tokens = "optional"` is an explicit, statically detectable opt-in to IMDSv1 compatibility.
  - Restricting the rule to launch templates and to the explicit `"optional"` value avoids guessing about account defaults, AMI defaults, or organization-wide controls.
  - The alternative candidates were either broader variants of the same idea or depended more heavily on deployment intent.
- Sources used:
  - https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template
  - https://github.com/hashicorp/terraform-provider-aws/issues/25909
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_launch_template_imdsv2_optional_tokens.go`
  - `rules/aws_launch_template_imdsv2_optional_tokens_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports launch templates explicitly configured with `metadata_options.http_tokens = "optional"`.
- Follow-up ideas:
  - Consider the same explicit-optional IMDS check for `aws_instance` if the repository wants parallel coverage outside launch templates.
  - Look for another low-noise rule around deprecated or invalid enum-style values in EC2 and RDS resources.

## 2026-03-23 - Cycle 3

- Goal: Extend low-noise IMDSv2 coverage to another common EC2 resource.
- Candidates investigated:
  - `awscx_instance_imdsv2_optional_tokens`
  - `awscx_db_instance_publicly_accessible`
  - `awscx_efs_file_system_missing_provisioned_throughput`
- Selected candidate:
  - `awscx_instance_imdsv2_optional_tokens`
- Why selected:
  - It reuses the same AWS guidance as the launch-template rule but covers direct `aws_instance` usage, which is still common in smaller modules and examples.
  - The explicit `"optional"` trigger remains low-noise and avoids guessing about account defaults or omitted configuration.
  - The alternative candidates were either more environment-dependent or needed deeper schema-specific validation logic.
- Sources used:
  - https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance
  - https://github.com/hashicorp/terraform-provider-aws/issues/10949

- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_instance_imdsv2_optional_tokens.go`
  - `rules/aws_instance_imdsv2_optional_tokens_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports `aws_instance` resources explicitly configured with `metadata_options.http_tokens = "optional"`.
- Follow-up ideas:
  - Look for a non-IMDS candidate next to keep the ruleset balanced across AWS services.
  - Explore simple enum or deprecated-argument checks in RDS, EFS, and EC2 resources.

## 2026-03-23 - Cycle 4

- Goal: Add a non-IMDS rule grounded in an AWS or provider requirement with low false-positive risk.
- Candidates investigated:
  - `awscx_efs_file_system_missing_provisioned_throughput`
  - `awscx_db_instance_publicly_accessible`
  - `awscx_lb_missing_deletion_protection`
- Selected candidate:
  - `awscx_efs_file_system_missing_provisioned_throughput`
- Why selected:
  - AWS explicitly requires a provisioned throughput value when an EFS file system uses `throughput_mode = "provisioned"`.
  - The condition is easy to detect statically and does not depend on account policy, workload intent, or organization-specific conventions.
  - The alternative candidates were more advisory and therefore more likely to be noisy in legitimate deployments.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system
  - https://docs.aws.amazon.com/efs/latest/ug/throughput-modes.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_efs_file_system_missing_provisioned_throughput.go`
  - `rules/aws_efs_file_system_missing_provisioned_throughput_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added an `ERROR` rule that reports `aws_efs_file_system` resources using provisioned throughput mode without `provisioned_throughput_in_mibps`.
- Follow-up ideas:
  - Look for another low-noise validity rule in EFS or RDS to keep the ruleset balanced across services.
  - Revisit advisory candidates only if they can be narrowed to explicit high-signal configurations.

## 2026-03-23 - Cycle 9

- Goal: Research another low-noise AWS/provider-specific rule and implement one practical warning.
- Candidates investigated:
  - `awscx_eks_addon_deprecated_resolve_conflicts`
  - `awscx_s3_bucket_deprecated_logging`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_eks_addon_deprecated_resolve_conflicts`
- Why selected:
  - The provider deprecation is explicit and actionable, and the replacement attributes map cleanly to AWS create/update behavior.
  - Detection is low-noise because the rule only flags explicit use of the deprecated attribute.
  - The alternatives were either another near-duplicate S3 deprecation warning or a security heuristic with materially higher intent sensitivity.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_addon
  - https://github.com/hashicorp/terraform-provider-aws/issues/27481
  - https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateAddon.html
  - https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateAddon.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_eks_addon_deprecated_resolve_conflicts.go`
  - `rules/aws_eks_addon_deprecated_resolve_conflicts_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated `resolve_conflicts` usage on `aws_eks_addon` and points users to the split create/update attributes.
- Follow-up ideas:
  - Revisit `aws_s3_bucket_deprecated_logging` if the repository wants another narrowly scoped provider-upgrade warning.
  - Look for a non-deprecation EKS rule with similarly explicit detection, such as an invalid or contradictory attribute combination.

## 2026-03-23 - Cycle 16

- Goal: Add another low-noise AWS validity rule grounded in an explicit provider/API constraint.
- Candidates investigated:
  - `awscx_ebs_volume_throughput_non_gp3`
  - `awscx_sqs_queue_invalid_fifo_throughput_limit`
  - `awscx_efs_file_system_provisioned_throughput_non_provisioned`
- Selected candidate:
  - `awscx_ebs_volume_throughput_non_gp3`
- Why selected:
  - The provider documentation and EC2 API both state that EBS `throughput` is valid only for `gp3`, making this a direct invalid-combination check rather than a best-practice heuristic.
  - Detection is simple and low-noise because it only reports explicit `throughput` usage with an omitted or non-`gp3` `type`.
  - The SQS candidate is still promising but slightly narrower in day-to-day usage, and the EFS inverse check overlaps heavily with the existing provisioned-throughput rule family.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume
  - https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_EbsBlockDevice.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_ebs_volume_throughput_non_gp3.go`
  - `rules/aws_ebs_volume_throughput_non_gp3_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added an `ERROR` rule that reports `aws_ebs_volume` resources using `throughput` without `type = "gp3"`.
- Follow-up ideas:
  - Revisit the SQS high-throughput FIFO attribute-combination rule for a future cycle.
  - Consider the inverse EFS throughput-mode rule only if it adds value beyond the existing missing-throughput check.

## 2026-03-23 - Cycle 10

- Goal: Add another low-noise AWS/provider validity rule with explicit RDS storage semantics.
- Candidates investigated:
  - `awscx_db_instance_storage_throughput_non_gp3`
  - `awscx_db_instance_dedicated_log_volume_without_provisioned_iops`
  - `awscx_ebs_volume_throughput_non_gp3`
- Selected candidate:
  - `awscx_db_instance_storage_throughput_non_gp3`
- Why selected:
  - The provider documentation explicitly states that `storage_throughput` can only be set when `storage_type = "gp3"`, making this a direct validity check rather than an organizational policy preference.
  - Detection is low-noise because it only needs explicit resource attributes and can conservatively skip unknown `storage_type` expressions.
  - The dedicated-log-volume candidate still needs tighter confirmation of the exact Terraform storage-type mapping for "Provisioned IOPS", and the EBS throughput idea was deprioritized in favor of the already-documented RDS backlog item.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
  - https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/db_instance.html.markdown
  - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Storage.html#gp3-storage
  - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_PIOPS.StorageTypes.html#USER_PIOPS.dlv
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_db_instance_storage_throughput_non_gp3.go`
  - `rules/aws_db_instance_storage_throughput_non_gp3_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added an `ERROR` rule that reports `aws_db_instance` resources setting `storage_throughput` without `storage_type = "gp3"`.
- Follow-up ideas:
  - Revisit `awscx_db_instance_dedicated_log_volume_without_provisioned_iops` after verifying the exact accepted `storage_type` set from provider or API references.
  - Explore another explicit storage-combination rule outside RDS to avoid clustering too many consecutive database checks.

## 2026-03-23 - Cycle 15

- Goal: Implement one more low-noise AWS provider rule grounded directly in current provider documentation.
- Candidates investigated:
  - `awscx_instance_deprecated_network_interface`
  - `awscx_db_instance_storage_throughput_non_gp3`
  - `awscx_db_instance_dedicated_log_volume_without_provisioned_iops`
- Selected candidate:
  - `awscx_instance_deprecated_network_interface`
- Why selected:
  - The `aws_instance` docs explicitly mark `network_interface` as deprecated and provide direct migration targets, so the rule is easy to explain and low risk.
  - Detection is simple because it only flags explicit deprecated block usage.
  - The RDS candidates are still promising, but one was deferred for a later cycle and the other was not selected to avoid stacking multiple RDS storage rules in a single pass.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance
  - https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/instance.html.markdown
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
  - https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/db_instance.html.markdown
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_instance_deprecated_network_interface.go`
  - `rules/aws_instance_deprecated_network_interface_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated `network_interface` blocks on `aws_instance` and points users to the documented replacement resources.
- Follow-up ideas:
  - Implement `awscx_db_instance_storage_throughput_non_gp3`.
  - Revisit `dedicated_log_volume` after tightening the accepted storage-type mapping.

## 2026-03-23 - Cycle 14

- Goal: Add another low-noise AWS validity rule from primary documentation.
- Candidates investigated:
  - `awscx_sqs_queue_fifo_name_suffix`
  - `awscx_sns_topic_fifo_name_suffix`
  - `awscx_db_instance_storage_throughput_requires_gp3`
- Selected candidate:
  - `awscx_sqs_queue_fifo_name_suffix`
- Why selected:
  - AWS documents the `.fifo` suffix as a hard requirement for FIFO SQS queue names, making this a direct validity check rather than an opinionated best-practice rule.
  - Detection is simple and low-noise because it only reports resources that explicitly set both `fifo_queue = true` and a non-compliant `name`.
  - The SNS variant is similarly valid but slightly less common in Terraform codebases, while the RDS storage-throughput candidate would need more careful handling around adjacent storage constraints.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue
  - https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_sqs_queue_fifo_name_suffix.go`
  - `rules/aws_sqs_queue_fifo_name_suffix_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `ERROR` rule that reports `aws_sqs_queue` resources with `fifo_queue = true` whose explicit `name` does not end with `.fifo`.
- Follow-up ideas:
  - Revisit `awscx_sns_topic_fifo_name_suffix` as a sibling low-noise rule for SNS FIFO topics.
  - Explore `aws_db_instance.storage_throughput` rules once the repository is ready to encode more provider-specific conditional constraints.

## 2026-03-23 - Cycle 10

- Goal: Continue the loop with one more explicit, low-noise AWS provider deprecation rule.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_replication_configuration`
  - `awscx_s3_bucket_deprecated_cors_rule`
  - `awscx_s3_bucket_deprecated_website`
  - `awscx_s3_bucket_deprecated_grant`
- Selected candidate:
  - `awscx_s3_bucket_deprecated_replication_configuration`
- Why selected:
  - The `aws_s3_bucket` documentation explicitly marks `replication_configuration` as deprecated and points to a single replacement resource.
  - Detection is low-noise because it only reports an explicitly configured inline block.
  - Compared with `cors_rule`, `website`, and `grant`, replication configuration is a higher-value setting to migrate cleanly while retaining the same simple implementation shape.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_replication_configuration
  - https://github.com/hashicorp/terraform-provider-aws/issues/20433
  - https://www.hashicorp.com/blog/terraform-aws-provider-4-0-refactors-s3-bucket-resource
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_s3_bucket_deprecated_replication_configuration.go`
  - `rules/aws_s3_bucket_deprecated_replication_configuration_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - pending
- Result:
  - Added a new `WARNING` rule that reports deprecated inline `replication_configuration` blocks on `aws_s3_bucket` and directs users to `aws_s3_bucket_replication_configuration`.
- Follow-up ideas:
  - Revisit `cors_rule`, `website`, and `grant` if the next cycle still favors low-risk migration checks.
  - Prefer a non-S3 validation rule in a later cycle if a similarly explicit AWS/provider requirement is available.

## 2026-03-23 - Cycle 11

- Goal: Continue the loop with another explicit S3 bucket deprecation rule that has a single migration target.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_website`
  - `awscx_s3_bucket_deprecated_cors_rule`
  - `awscx_s3_bucket_deprecated_grant`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_s3_bucket_deprecated_website`
- Why selected:
  - The `aws_s3_bucket` documentation explicitly marks the inline `website` block as deprecated and points users to `aws_s3_bucket_website_configuration`.
  - Detection is low-noise because it only reports an explicitly configured inline block.
  - Compared with `cors_rule` and `grant`, the replacement resource is straightforward and the static website use case is common enough to justify a dedicated migration warning.
- Sources used:
  - https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteHosting.html
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration
  - https://github.com/hashicorp/terraform-provider-aws/issues/20433
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_s3_bucket_deprecated_website.go`
  - `rules/aws_s3_bucket_deprecated_website_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - pending
- Result:
  - Added a new `WARNING` rule that reports deprecated inline `website` blocks on `aws_s3_bucket` and directs users to `aws_s3_bucket_website_configuration`.
- Follow-up ideas:
  - Revisit `cors_rule` and `grant` if the next cycle still favors low-risk S3 migration checks.
  - Prefer a non-S3 candidate next if a comparably explicit AWS/provider validation rule can be identified.

## 2026-03-23 - Cycle 13

- Goal: Continue the loop with a non-S3, low-noise AWS/provider-specific rule.
- Candidates investigated:
  - `awscx_api_gateway_deployment_deprecated_stage_management`
  - `awscx_eip_deprecated_vpc`
  - `awscx_lb_listener_invalid_mutual_authentication`
- Selected candidate:
  - `awscx_api_gateway_deployment_deprecated_stage_management`
- Why selected:
  - The AWS provider explicitly deprecated `stage_name`, `stage_description`, and `canary_settings` on `aws_api_gateway_deployment`, with a clear migration path to `aws_api_gateway_stage`.
  - AWS documentation models deployments and stages as separate resources, so the rule aligns with both provider direction and service concepts.
  - The alternatives were weaker for this repository cycle: `aws_eip.vpc` is another narrow deprecation but lower practical value, and ALB mutual-auth validation would require more nuanced static semantics with higher false-positive risk.
- Sources used:
  - https://github.com/hashicorp/terraform-provider-aws/issues/39957
  - https://github.com/hashicorp/terraform-provider-aws/issues/39958
  - https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-deployments.html
  - https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-stages.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_api_gateway_deployment_deprecated_stage_management.go`
  - `rules/aws_api_gateway_deployment_deprecated_stage_management_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated stage-management fields on `aws_api_gateway_deployment` and directs users to `aws_api_gateway_stage`.
- Follow-up ideas:
  - Revisit `awscx_eip_deprecated_vpc` if the repository wants another small migration warning next.
  - Look for a direct API Gateway validity rule after this deprecation cleanup, but keep detection limited to explicit contradictory configuration.

## 2026-03-23 - Cycle 12

- Goal: Continue the implementation loop with a high-value, low-noise AWS provider migration rule.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_lifecycle_rule`
  - `awscx_s3_bucket_deprecated_replication_configuration`
  - `awscx_db_instance_publicly_accessible`
  - `awscx_lb_deletion_protection_disabled`
- Selected candidate:
  - `awscx_s3_bucket_deprecated_lifecycle_rule`
- Why selected:
  - The AWS provider explicitly split inline S3 bucket lifecycle management into `aws_s3_bucket_lifecycle_configuration`, making the migration path concrete and low-ambiguity.
  - Detection is low-noise because the rule only reports explicit inline `lifecycle_rule` blocks and does not infer broader bucket policy intent.
  - The alternative S3 replication candidate is similarly valid but slightly less common, while the RDS and ALB candidates remain more advisory and environment-sensitive.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration
  - https://github.com/hashicorp/terraform-provider-aws/issues/20433
  - https://www.hashicorp.com/blog/terraform-aws-provider-4-0-refactors-s3-bucket-resource
  - https://docs.aws.amazon.com/securityhub/latest/userguide/exposure-rds.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_s3_bucket_deprecated_lifecycle_rule.go`
  - `rules/aws_s3_bucket_deprecated_lifecycle_rule_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated inline `lifecycle_rule` usage on `aws_s3_bucket` and directs users to `aws_s3_bucket_lifecycle_configuration`.
- Follow-up ideas:
  - Continue the S3 split-resource migration family with `replication_configuration` if the repository wants another near-zero-noise provider deprecation warning.
  - Prefer a non-S3 candidate in the next cycle if one can match the same low false-positive bar.

## 2026-03-23 - Cycle 11

- Goal: Continue the loop with one more low-noise AWS/provider-specific rule and keep the change set small.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_logging`
  - `awscx_db_instance_publicly_accessible`
  - `awscx_lb_missing_deletion_protection`
- Selected candidate:
  - `awscx_s3_bucket_deprecated_logging`
- Why selected:
  - The provider's S3 bucket refactor plan explicitly lists inline `logging` among the arguments deprecated in favor of standalone resources.
  - Detection stays low-noise because the rule only flags an explicit inline `logging` block on `aws_s3_bucket`.
  - The RDS and load balancer alternatives remain more policy-driven and therefore riskier in repositories that intentionally expose public services.
- Sources used:
  - https://developer.hashicorp.com/validated-patterns/terraform/upgrade-terraform-provider
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_logging
  - https://github.com/hashicorp/terraform-provider-aws/issues/20433
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_s3_bucket_deprecated_logging.go`
  - `rules/aws_s3_bucket_deprecated_logging_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated inline `logging` blocks on `aws_s3_bucket` and directs users to `aws_s3_bucket_logging`.
- Follow-up ideas:
  - Consider `aws_s3_bucket_deprecated_lifecycle_rule` if the repository wants to complete more of the S3 split-resource migration family.
  - Prefer a non-S3 candidate next to keep AWS service coverage balanced.

## 2026-03-23 - Cycle 10

- Goal: Continue the rule loop with another low-noise AWS/provider-specific deprecation check outside the recent S3 and EKS work.
- Candidates investigated:
  - `awscx_launch_template_deprecated_elastic_gpu_specifications`
  - `awscx_s3_bucket_deprecated_logging`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_launch_template_deprecated_elastic_gpu_specifications`
- Why selected:
  - AWS documents Elastic Graphics as deprecated with an end-of-life date, and the provider has an open issue to deprecate and remove `elastic_gpu_specifications`.
  - Detection is low-noise because the rule only flags an explicit deprecated block on `aws_launch_template`.
  - The alternative S3 logging candidate remains viable but is another near-duplicate S3 split-resource warning, while the RDS public-access candidate is still materially more policy-driven.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template
  - https://github.com/hashicorp/terraform-provider-aws/issues/37589
  - https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_ElasticGpuSpecificationResponse.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_launch_template_deprecated_elastic_gpu_specifications.go`
  - `rules/aws_launch_template_deprecated_elastic_gpu_specifications_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - pending
- Result:
  - Added a new `WARNING` rule that reports deprecated `elastic_gpu_specifications` usage on `aws_launch_template`.
- Follow-up ideas:
  - Consider the sibling `aws_instance` variant only if the provider documentation clearly marks that path as deprecated too.
  - Revisit `aws_s3_bucket_deprecated_logging` when another small deprecation-only change set is acceptable.

## 2026-03-23 - Cycle 5

- Goal: Continue adding non-advisory AWS validity rules with simple static detection.
- Candidates investigated:
  - `awscx_ebs_volume_missing_iops`
  - `awscx_db_instance_publicly_accessible`
  - `awscx_lb_missing_deletion_protection`
- Selected candidate:
  - `awscx_ebs_volume_missing_iops`
- Why selected:
  - AWS directly requires `iops` for EBS volume types `io1` and `io2`.
  - The detection is explicit and does not depend on organization policy or runtime context.
  - The alternative candidates were more policy-oriented and therefore less suitable for an always-on low-noise rule.
- Sources used:
  - https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_CreateVolume.html
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_ebs_volume_missing_iops.go`
  - `rules/aws_ebs_volume_missing_iops_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added an `ERROR` rule that reports `aws_ebs_volume` resources using `io1` or `io2` without `iops`.
- Follow-up ideas:
  - Explore another EBS or RDS rule where a required companion attribute is missing for a specific enum value.
  - Revisit deletion protection only if the repository decides to add more policy-style warnings.

## 2026-03-23 - Cycle 6

- Goal: Extend low-noise required-companion-attribute checks into RDS.
- Candidates investigated:
  - `awscx_db_instance_missing_iops`
  - `awscx_lb_missing_deletion_protection`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_db_instance_missing_iops`
- Why selected:
  - The RDS API explicitly requires IOPS for storage types `io1`, `io2`, and `gp3`.
  - The detection is straightforward and based on a provider-facing requirement instead of environment policy.
  - The alternative candidates were more advisory and therefore less suitable as always-on checks.
- Sources used:
  - https://docs.aws.amazon.com/goto/WebAPI/rds-2014-10-31/CreateDBInstance
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_db_instance_missing_iops.go`
  - `rules/aws_db_instance_missing_iops_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added an `ERROR` rule that reports `aws_db_instance` resources using `io1`, `io2`, or `gp3` storage without `iops`.
- Follow-up ideas:
  - Look for the next low-noise RDS or EC2 rule involving another required companion attribute.
  - Keep deletion-protection-style warnings deferred unless they can be narrowed to very explicit unsafe configurations.

## 2026-03-23 - Cycle 7

- Goal: Add another low-noise AWS-specific rule without reusing the recent required-attribute pattern again.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_versioning`
  - `awscx_s3_bucket_deprecated_server_side_encryption_configuration`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_s3_bucket_deprecated_versioning`
- Why selected:
  - HashiCorp's AWS provider upgrade guidance explicitly moves bucket versioning management from inline `aws_s3_bucket.versioning` to `aws_s3_bucket_versioning`.
  - Detection is simple and low-noise because it only reports explicit inline usage and does not infer security posture or module intent.
  - The alternative S3 encryption candidate was almost identical in shape, so it was deferred to keep this change set focused; the RDS public-access candidate remained more policy-driven and therefore noisier.
- Sources used:
  - https://developer.hashicorp.com/validated-patterns/terraform/upgrade-terraform-provider
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_versioning
  - https://github.com/hashicorp/terraform-provider-aws/issues/20433
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_s3_bucket_deprecated_versioning.go`
  - `rules/aws_s3_bucket_deprecated_versioning_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated inline `versioning` blocks on `aws_s3_bucket` and directs users to `aws_s3_bucket_versioning`.
- Follow-up ideas:
  - Implement the sibling S3 deprecation rule for inline `server_side_encryption_configuration` if the repository wants another low-risk migration check next.
  - Continue balancing deprecation rules with direct AWS validity checks in other services.

## 2026-03-23 - Cycle 8

- Goal: Continue the rule loop with another low-noise AWS-specific deprecation check.
- Candidates investigated:
  - `awscx_s3_bucket_deprecated_server_side_encryption_configuration`
  - `awscx_s3_bucket_deprecated_logging`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_s3_bucket_deprecated_server_side_encryption_configuration`
- Why selected:
  - HashiCorp's upgrade guidance explicitly shows `aws_s3_bucket.server_side_encryption_configuration` producing a deprecation warning and migrating to `aws_s3_bucket_server_side_encryption_configuration`.
  - Detection is simple and low-noise because it only reports explicit inline block usage on `aws_s3_bucket`.
  - The inline logging candidate is valid but less prominent in the migration guidance, while the RDS public-access candidate remains more policy-driven than provider-driven.
- Sources used:
  - https://developer.hashicorp.com/validated-patterns/terraform/upgrade-terraform-provider
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_server_side_encryption_configuration
  - https://github.com/hashicorp/terraform-provider-aws/issues/20433
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_s3_bucket_deprecated_server_side_encryption_configuration.go`
  - `rules/aws_s3_bucket_deprecated_server_side_encryption_configuration_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated inline `server_side_encryption_configuration` blocks on `aws_s3_bucket` and directs users to `aws_s3_bucket_server_side_encryption_configuration`.
- Follow-up ideas:
  - Revisit `logging` or `lifecycle_rule` if the repository wants to continue the S3 bucket split-resource migration family.
  - Prefer a non-S3 candidate in the following cycle to keep service coverage broad.

## 2026-03-23 - Cycle 9

- Goal: Research another low-noise AWS/provider-specific rule and implement one practical warning.
- Candidates investigated:
  - `awscx_eks_addon_deprecated_resolve_conflicts`
  - `awscx_s3_bucket_deprecated_logging`
  - `awscx_db_instance_publicly_accessible`
- Selected candidate:
  - `awscx_eks_addon_deprecated_resolve_conflicts`
- Why selected:
  - The provider deprecation is explicit and actionable, and the replacement attributes map cleanly to AWS create/update behavior.
  - Detection is low-noise because the rule only flags explicit use of the deprecated attribute.
  - The alternatives were either another near-duplicate S3 deprecation warning or a security heuristic with materially higher intent sensitivity.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_addon
  - https://github.com/hashicorp/terraform-provider-aws/issues/27481
  - https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateAddon.html
  - https://docs.aws.amazon.com/eks/latest/APIReference/API_UpdateAddon.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_eks_addon_deprecated_resolve_conflicts.go`
  - `rules/aws_eks_addon_deprecated_resolve_conflicts_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `WARNING` rule that reports deprecated `resolve_conflicts` usage on `aws_eks_addon` and points users to the split create/update attributes.
- Follow-up ideas:
  - Revisit `aws_s3_bucket_deprecated_logging` if the repository wants another narrowly scoped provider-upgrade warning.
  - Look for a non-deprecation EKS rule with similarly explicit detection, such as an invalid or contradictory attribute combination.

## 2026-03-23 - Cycle 17

- Goal: Add another low-noise AWS-specific rule with a clear provider-documented invalid attribute combination.
- Candidates investigated:
  - `awscx_db_instance_dedicated_log_volume_non_io1_io2`
  - `awscx_efs_file_system_provisioned_throughput_non_provisioned_mode`
  - `awscx_ebs_volume_multi_attach_non_io1_io2`
- Selected candidate:
  - `awscx_db_instance_dedicated_log_volume_non_io1_io2`
- Why selected:
  - The provider documentation explicitly states that `dedicated_log_volume` requires Provisioned IOPS storage, and Terraform `storage_type` maps that requirement directly to `io1` or `io2`.
  - Detection is straightforward and low-noise because it only reports explicit `dedicated_log_volume = true` with an omitted or incompatible `storage_type`.
  - The EFS candidate is also valid but is a close sibling of an existing EFS rule, while the EBS Multi-Attach candidate needs a little more care because Windows support differs between `io1` and `io2`.
- Sources used:
  - https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance
  - https://raw.githubusercontent.com/hashicorp/terraform-provider-aws/main/website/docs/r/db_instance.html.markdown
  - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_PIOPS.StorageTypes.html
- Files changed:
  - `main.go`
  - `README.md`
  - `rules/aws_db_instance_dedicated_log_volume_non_io1_io2.go`
  - `rules/aws_db_instance_dedicated_log_volume_non_io1_io2_test.go`
  - `notes/rule-backlog.md`
  - `notes/research-log.md`
- Tests run:
  - `go test ./...`
- Result:
  - Added a new `ERROR` rule that reports `aws_db_instance` resources using `dedicated_log_volume = true` without an explicit `io1` or `io2` storage type.
- Follow-up ideas:
  - Implement the converse EFS throughput rule if the repository wants another explicit invalid-combination check next.
  - Revisit the EBS Multi-Attach candidate with OS-sensitive wording or tighter Terraform-side scope.
