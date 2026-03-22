# AGENTS.md

This file defines how agents should work in this repository.

## Goal

This repository builds a custom TFLint ruleset focused on AWS in Terraform usage.

The primary ongoing task is:

1. Discover useful lint rule ideas from [Terraform Registry provider documentation for AWS resources](https://registry.terraform.io/providers/hashicorp/aws/latest/docs), [`terraform-provider-aws`](https://github.com/hashicorp/terraform-provider-aws) issues/PRs, and [AWS documentation](https://docs.aws.amazon.com/).
2. Select one high-value, low-noise rule candidate.
3. Implement the rule in this repository.
4. Add or update tests.
5. Run `go test ./...`.
7. Update user-facing documentation and research notes.
8. Prepare a clean, reviewable change set. Run `git commit -m <message>` and `git push origin main`.

Unless the user says otherwise, optimize for practical rules that catch common mistakes, outdated usage, unsafe defaults, or provider-specific pitfalls with low false-positive risk.

## Repository Context

This repository contains example rules. These can be used as examples of how to use the TFLint Plugin SDK.
For other usage, please refer to the [GoDoc for tflint-plugin-sdk](https://pkg.go.dev/github.com/terraform-linters/tflint-plugin-sdk).

## Source Priority

When researching new rule ideas, use sources in this order:

1. [Terraform Registry provider documentation for AWS resources](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
2. [`terraform-provider-aws`](https://github.com/hashicorp/terraform-provider-aws) issues and pull requests
3. [AWS official documentation](https://docs.aws.amazon.com/)
4. Existing provider behavior visible from resource schema or common configuration patterns

Use official or primary sources whenever possible. Record the source URLs used for the chosen rule.

## Rule Selection Criteria

Prefer rule ideas that score well on all of the following:

1. High practical value
2. Clear AWS or provider-specific rationale
3. Low false-positive risk
4. Feasible static detection from Terraform configuration
5. Easy to explain in a short lint message

Good candidates usually have one or more of these properties:

1. They reflect an AWS recommendation or requirement.
2. They prevent a common misconfiguration discussed in provider issues.
3. They detect deprecated or superseded configuration patterns.
4. They catch insecure or surprising defaults.

Avoid or deprioritize candidates that:

1. Require runtime or account context.
2. Depend on values that are usually unknown until apply.
3. Overlap heavily with existing well-known checks unless this repository adds clear AWS-specific value.
4. Have ambiguous semantics across modules or organizational conventions.

## Default Work Cycle

Unless the user gives a different instruction, one work cycle means:

1. Research 2 to 4 candidate rules.
2. Compare them briefly.
3. Pick the best candidate.
4. Implement exactly one rule.
5. Add tests.
6. Update documentation/backlog notes.
7. Run verification and record the result.

If no good candidate is found, do not force an implementation. Instead:

1. Document the investigated candidates.
2. Explain why they were rejected or deferred.

DO NOT interrupt these work cycles until instructed to do so by the user. DO NOT ask the user for confirmation.

## Expected Outputs Per Implemented Rule

For each accepted rule, aim to update all relevant areas:

1. Rule implementation under `rules/`
2. Tests under `rules/*_test.go`
3. Rule registration in `main.go`
4. Rule listing and description in `README.md`
5. Research/backlog notes in `notes/*`

## Implementation Guidance

When implementing a rule:

1. Use a clear AWS-specific rule name.
2. Set severity intentionally; default to `ERROR` only when the configuration is clearly invalid or risky.
3. Fill in `Link()` with the most relevant documentation or issue URL when possible.
4. Keep issue messages specific and actionable.
5. Prefer checking explicit configuration over inferred intent.

Design rules so that they are understandable from code review alone. Small, direct checks are preferred over clever but fragile logic.

## Testing Guidance

Each real rule should normally include tests for:

1. A configuration that should trigger an issue
2. A configuration that should not trigger an issue
3. Edge cases that are important for false-positive control

When expression values may be unknown or variable-driven, be conservative. Prefer skipping ambiguous cases over reporting speculative issues.

## Documentation Guidance

When updating `README.md`, keep the rule table accurate. For each real rule, include:

1. Rule name
2. Short description
3. Severity
4. Enabled status
5. Reference link when available

When a rule is implemented or updated:

1. update `notes/rule-backlog.md`
2. append a short entry to `notes/research-log.md`
3. keep the source URLs used for the rule

## How To Report Progress

After each substantial work cycle, report:

1. Candidates investigated
2. Sources used
3. Selected rule and why
4. Files changed
5. Tests run and results
6. Remaining risks or follow-up ideas

Keep the report concise, but include enough detail for the next cycle to continue smoothly.

## Safety And Change Discipline

Do not revert unrelated user changes.

Prefer small, reviewable changes. If a candidate rule would require broad refactoring, either:

1. Split the refactor from the rule implementation, or
2. Defer the rule and choose a smaller candidate first

## Biases For This Repository

In ambiguous cases, prefer:

1. Lower false-positive risk over broader coverage
2. Rules grounded in AWS guidance over purely stylistic preferences
3. Incremental repository improvement over large template cleanup
4. Replacing example value with real value over adding more examples
