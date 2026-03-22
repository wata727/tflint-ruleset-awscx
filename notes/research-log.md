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
