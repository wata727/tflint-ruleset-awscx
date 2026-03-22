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
|awscx_api_gateway_deployment_deprecated_stage_management|Warn on deprecated stage management fields on `aws_api_gateway_deployment`|WARNING|✔|https://github.com/hashicorp/terraform-provider-aws/issues/39957|
|awscx_db_instance_dedicated_log_volume_non_io1_io2|Disallow `dedicated_log_volume` unless `storage_type` is `io1` or `io2`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_db_instance_missing_iops|Require `iops` when `storage_type` is `io1`, `io2`, or `gp3`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_db_instance_storage_throughput_non_gp3|Disallow `storage_throughput` unless `storage_type` is `gp3`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_ebs_volume_missing_iops|Require `iops` when `type` is `io1` or `io2`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume|
|awscx_ebs_volume_throughput_non_gp3|Disallow `throughput` unless `type` is `gp3`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume|
|awscx_dynamodb_table_invalid_stream_view_type|Disallow invalid `stream_view_type`|ERROR|✔|https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_StreamSpecification.html|
|awscx_eks_addon_deprecated_resolve_conflicts|Warn on deprecated `resolve_conflicts` on `aws_eks_addon`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_addon|
|awscx_efs_file_system_missing_provisioned_throughput|Require `provisioned_throughput_in_mibps` when `throughput_mode` is `provisioned`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system|
|awscx_instance_deprecated_network_interface|Warn on deprecated `network_interface` blocks on `aws_instance`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance|
|awscx_instance_imdsv2_optional_tokens|Warn when `metadata_options.http_tokens` on `aws_instance` explicitly allows IMDSv1|WARNING|✔|https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html|
|awscx_launch_template_deprecated_elastic_gpu_specifications|Warn on deprecated `elastic_gpu_specifications` on `aws_launch_template`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template|
|awscx_launch_template_imdsv2_optional_tokens|Warn when `metadata_options.http_tokens` explicitly allows IMDSv1|WARNING|✔|https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html|
|awscx_s3_bucket_deprecated_acl|Warn on deprecated inline `acl` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_lifecycle_rule|Warn on deprecated inline `lifecycle_rule` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_logging|Warn on deprecated inline `logging` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_replication_configuration|Warn on deprecated inline `replication_configuration` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_server_side_encryption_configuration|Warn on deprecated inline `server_side_encryption_configuration` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_versioning|Warn on deprecated inline `versioning` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_website|Warn on deprecated inline `website` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_security_group_invalid_protocol|Disallow invalid `protocol`|ERROR|✔|https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html|
|awscx_sqs_queue_fifo_name_suffix|Require FIFO `aws_sqs_queue.name` values to end with `.fifo`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue|
