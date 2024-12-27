resource "google_dns_managed_zone" "tomedome" {
  name        = "tomedome-public"
  dns_name    = "tomedome.io."
  description = "Tomedome public DNS zone"
}

resource "google_dns_record_set" "frontend" {
  name = "dota.${google_dns_managed_zone.tomedome.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = google_dns_managed_zone.tomedome.name

  rrdatas = [
    google_compute_global_forwarding_rule.tomedome.ip_address
  ]
}

resource "google_dns_record_set" "api" {
  name = "api.${google_dns_managed_zone.tomedome.dns_name}"
  type = "A"
  ttl  = 300

  managed_zone = google_dns_managed_zone.tomedome.name

  rrdatas = [
    google_compute_global_forwarding_rule.tomedome.ip_address
  ]
}
