variable "zone_name" {
  type        = string
  description = "name of the zone"
}

variable "records" {
  type = list(
    object({
      name    = string
      type    = string
      records = list(string)
      ttl     = number
    })
  )
  description = "List of record object belonging to this zone"
}
