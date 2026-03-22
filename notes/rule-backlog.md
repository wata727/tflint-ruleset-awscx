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

- Status: deferred
- Resource(s): `aws_s3_bucket`
- Short description: Warn when deprecated inline `logging` is used on `aws_s3_bucket`.
- Why it matters: This is another S3 bucket argument included in the provider's split-resource deprecation plan.
- Detection approach: Flag any `logging` block present on `aws_s3_bucket`.
- False-positive risk: low
- Implementation difficulty: low
- Overlap notes: Similar to the selected S3 deprecation rules; deferred to avoid shipping too many nearly identical warnings in one pass.
- Selected on:
- Implemented on:

### Sources

- HashiCorp docs: https://developer.hashicorp.com/validated-patterns/terraform/upgrade-terraform-provider
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket
- Terraform Registry docs: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_logging
- terraform-provider-aws issue/PR: https://github.com/hashicorp/terraform-provider-aws/issues/20433

### Notes

- Viable follow-up if the repository keeps extending the S3 bucket deprecation family.

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
