# Route53 Zone Module

<!-- BEGIN_TF_DOCS -->


## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

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
<!-- END_TF_DOCS -->
