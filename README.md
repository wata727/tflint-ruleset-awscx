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
|awscx_autoscaling_group_invalid_max_instance_lifetime|Require `max_instance_lifetime` to be `0` or between `86400` and `31536000` seconds|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/autoscaling_group|
|awscx_cloudfront_distribution_minimum_protocol_version_default_certificate|Disallow `viewer_certificate.minimum_protocol_version` when `viewer_certificate.cloudfront_default_certificate = true`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudfront_distribution|
|awscx_cloudwatch_log_group_delivery_retention_in_days|Warn when `retention_in_days` is set on `aws_cloudwatch_log_group` with `log_group_class = "DELIVERY"`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group|
|awscx_db_instance_blue_green_update_without_backup_retention|Require `backup_retention_period > 0` when `blue_green_update.enabled = true`|ERROR|✔|https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/blue-green-deployments-creating.html|
|awscx_db_instance_database_insights_advanced_requirements|Require Performance Insights and retention >= 465 when `database_insights_mode = "advanced"`|ERROR|✔|https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_DatabaseInsights.TurningOnAdvanced.html|
|awscx_db_instance_dedicated_log_volume_non_io1_io2|Disallow `dedicated_log_volume` unless `storage_type` is `io1` or `io2`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_db_instance_enhanced_monitoring_role_requirements|Require `monitoring_role_arn` when `monitoring_interval` enables Enhanced Monitoring and disallow it for `0`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_db_instance_manage_master_user_password_conflict|Disallow `password` and `password_wo` when `manage_master_user_password = true`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_db_instance_performance_insights_arguments_without_enabled|Disallow Performance Insights KMS or retention arguments unless `performance_insights_enabled = true`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_db_instance_missing_iops|Require `iops` when `storage_type` is `io1`, `io2`, or `gp3`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_db_instance_storage_throughput_non_gp3|Disallow `storage_throughput` unless `storage_type` is `gp3`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/db_instance|
|awscx_ebs_volume_missing_iops|Require `iops` when `type` is `io1` or `io2`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume|
|awscx_ebs_volume_throughput_non_gp3|Disallow `throughput` unless `type` is `gp3`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume|
|awscx_eip_instance_network_interface_conflict|Disallow setting both `instance` and `network_interface` on `aws_eip`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip|
|awscx_ecs_service_deployment_maximum_percent_daemon|Disallow `deployment_maximum_percent` when `aws_ecs_service.scheduling_strategy = "DAEMON"`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service|
|awscx_ecs_service_health_check_grace_period_without_load_balancer|Disallow `health_check_grace_period_seconds` unless `aws_ecs_service` has a `load_balancer` block|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecs_service|
|awscx_dynamodb_table_invalid_stream_view_type|Disallow invalid `stream_view_type`|ERROR|✔|https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_StreamSpecification.html|
|awscx_eks_addon_deprecated_resolve_conflicts|Warn on deprecated `resolve_conflicts` on `aws_eks_addon`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eks_addon|
|awscx_efs_file_system_missing_provisioned_throughput|Require `provisioned_throughput_in_mibps` when `throughput_mode` is `provisioned`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system|
|awscx_efs_file_system_kms_key_without_encrypted|Disallow `kms_key_id` unless `encrypted = true` on `aws_efs_file_system`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system|
|awscx_efs_file_system_provisioned_throughput_non_provisioned|Disallow `provisioned_throughput_in_mibps` unless `throughput_mode` is `provisioned`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/efs_file_system|
|awscx_instance_deprecated_network_interface|Warn on deprecated `network_interface` blocks on `aws_instance`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance|
|awscx_instance_imdsv2_optional_tokens|Warn when `metadata_options.http_tokens` on `aws_instance` explicitly allows IMDSv1|WARNING|✔|https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html|
|awscx_lb_listener_alpn_policy_non_tls|Disallow `alpn_policy` unless `aws_lb_listener.protocol` is `TLS`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener|
|awscx_lb_listener_missing_certificate_arn|Require `certificate_arn` when `aws_lb_listener.protocol` is `HTTPS` or `TLS`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener|
|awscx_lb_listener_mutual_authentication_verify_requirements|Require `trust_store_arn` for `mutual_authentication.mode = "verify"` and disallow verify-only attributes otherwise|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener|
|awscx_lb_listener_missing_ssl_policy|Require `ssl_policy` when `aws_lb_listener.protocol` is `HTTPS` or `TLS`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_listener|
|awscx_lb_target_group_lambda_top_level_attributes|Disallow `port`, `protocol`, and `vpc_id` when `aws_lb_target_group.target_type` is `lambda`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group|
|awscx_lb_target_group_matcher_non_http_health_check|Disallow `health_check.matcher` unless `health_check.protocol` is `HTTP` or `HTTPS`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group|
|awscx_lb_target_group_protocol_version_non_http|Disallow `protocol_version` unless `aws_lb_target_group.protocol` is `HTTP` or `HTTPS`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/lb_target_group|
|awscx_launch_template_deprecated_elastic_inference_accelerator|Warn on deprecated `elastic_inference_accelerator` on `aws_launch_template`|WARNING|✔|https://github.com/hashicorp/terraform-provider-aws/issues/41101|
|awscx_launch_template_deprecated_elastic_gpu_specifications|Warn on deprecated `elastic_gpu_specifications` on `aws_launch_template`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_template|
|awscx_launch_template_imdsv2_optional_tokens|Warn when `metadata_options.http_tokens` explicitly allows IMDSv1|WARNING|✔|https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-IMDS-new-instances.html|
|awscx_s3_bucket_deprecated_acceleration_status|Warn on deprecated inline `acceleration_status` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_acl|Warn on deprecated inline `acl` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_configuration_expected_bucket_owner_deprecated|Warn on deprecated `expected_bucket_owner` across S3 bucket configuration sub-resources|WARNING|✔|https://github.com/hashicorp/terraform-provider-aws/pull/46262|
|awscx_s3_bucket_deprecated_grant|Warn on deprecated inline `grant` blocks on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_lifecycle_rule|Warn on deprecated inline `lifecycle_rule` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_logging|Warn on deprecated inline `logging` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_object_lock_configuration|Warn on deprecated inline `object_lock_configuration` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_policy|Warn on deprecated inline `policy` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_request_payer|Warn on deprecated inline `request_payer` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_replication_configuration|Warn on deprecated inline `replication_configuration` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_server_side_encryption_configuration|Warn on deprecated inline `server_side_encryption_configuration` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_versioning|Warn on deprecated inline `versioning` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_s3_bucket_deprecated_website|Warn on deprecated inline `website` on `aws_s3_bucket`|WARNING|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket|
|awscx_sns_topic_fifo_attributes_without_fifo_topic|Disallow FIFO-only `aws_sns_topic` attributes unless `fifo_topic = true`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic|
|awscx_sns_topic_fifo_name_suffix|Require FIFO `aws_sns_topic.name` values to end with `.fifo`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sns_topic|
|awscx_security_group_all_protocol_nonzero_ports|Require `from_port = 0` and `to_port = 0` when `protocol = "-1"` on inline security group rules|ERROR|✔|https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html|
|awscx_security_group_invalid_protocol|Disallow invalid `protocol`|ERROR|✔|https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html|
|awscx_sfn_state_machine_log_destination_missing_wildcard|Require `logging_configuration.log_destination` on `aws_sfn_state_machine` to end with `:*`|ERROR|✔|https://docs.aws.amazon.com/step-functions/latest/apireference/API_CloudWatchLogsLogGroup.html|
|awscx_sqs_queue_fifo_name_suffix|Require FIFO `aws_sqs_queue.name` values to end with `.fifo`|ERROR|✔|https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue|
