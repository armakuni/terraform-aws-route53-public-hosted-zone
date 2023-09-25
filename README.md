# AWS Route53 Public Hosted Zone Terraform Module

Terraform module to easily provision a public hosted zone with given records on Route53.

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_route53_record.record](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_record) | resource |
| [aws_route53_zone.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_records"></a> [records](#input\_records) | n/a | <pre>list(<br>    object({<br>      name    = string<br>      type    = string<br>      records = list(string)<br>      ttl     = number<br>    })<br>  )</pre> | n/a | yes |
| <a name="input_zone_name"></a> [zone\_name](#input\_zone\_name) | n/a | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_zone_id"></a> [zone\_id](#output\_zone\_id) | n/a |
<!-- END_TF_DOCS -->
