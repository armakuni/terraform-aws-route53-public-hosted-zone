variable "zone_name" {
  type = string
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
  validation {
    condition     = alltrue([for record in var.records : contains(["A", "CNAME", "MX", "NS", "TXT", "SOA", "SPF"], record.type)])
    error_message = "Only valid types permitted (A, CNAME, MX, NS, TXT, SOA, SPF)"
  }
}
