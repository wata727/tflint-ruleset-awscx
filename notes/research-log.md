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
  - pending
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
  - pending
- Result:
  - Added a new `WARNING` rule that reports deprecated `resolve_conflicts` usage on `aws_eks_addon` and points users to the split create/update attributes.
- Follow-up ideas:
  - Revisit `aws_s3_bucket_deprecated_logging` if the repository wants another narrowly scoped provider-upgrade warning.
  - Look for a non-deprecation EKS rule with similarly explicit detection, such as an invalid or contradictory attribute combination.

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
