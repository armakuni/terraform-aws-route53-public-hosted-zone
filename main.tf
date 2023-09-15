resource "aws_route53_zone" "this" {
  name = var.zone_name
}

resource "aws_route53_record" "record" {
  for_each = { for record in var.records : "name=${record.name},type=${record.type}" => record }

  zone_id = aws_route53_zone.this.zone_id
  name    = each.value.name == "" ? var.zone_name : "${each.value.name}.${var.zone_name}"
  type    = each.value.type
  records = each.value.records
  ttl     = each.value.ttl
}
