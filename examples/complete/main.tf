module "test_route53_zone" {
  source    = "../.."
  zone_name = var.zone_name
  records   = var.records
}
