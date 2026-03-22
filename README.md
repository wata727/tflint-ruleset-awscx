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

## Local verification

For the repository-level verification flow, use:

```shell
go test ./...
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |
|awscx_dynamodb_table_invalid_stream_view_type|Disallow invalid `stream_view_type`|ERROR|✔|https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_StreamSpecification.html|
|awscx_instance_imdsv2_optional_tokens|Warn when `metadata_options.http_tokens` on `aws_instance` explicitly allows IMDSv1|WARNING|✔|https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html|
|awscx_launch_template_imdsv2_optional_tokens|Warn when `metadata_options.http_tokens` explicitly allows IMDSv1|WARNING|✔|https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html|
|awscx_s3_bucket_deprecated_acl|Warn on deprecated inline `acl` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_security_group_invalid_protocol|Disallow invalid `protocol`|ERROR|✔|https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html|
