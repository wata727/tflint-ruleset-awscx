# TFLint ruleset for terraform-provider-aws powered by Codex

TFLint ruleset plugin for Terraform AWS Provider powered by Codex.

## Requirements

- TFLint v0.46+
- Go v1.26

## Installation

Run the following command:

```shell
make install
```

Test the installed plugin behavior:

```shell
tflint
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|awscx_dynamodb_table_invalid_stream_view_type|Disallow invalid `stream_view_type`|ERROR|✔||
|awscx_security_group_invalid_protocol|Disallow invalid `protocol`|ERROR|✔||
