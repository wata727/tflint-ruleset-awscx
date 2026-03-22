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
