resource "google_compute_url_map" "tomedome" {

  name        = "tomedome-static"
  description = "Tomedome static site URL map"

  default_service = google_compute_backend_bucket.static_site_backend.self_link

  host_rule {
    hosts        = ["api.tomedome.io"]
    path_matcher = "api"
  }

  path_matcher {
    name            = "api"
    default_service = google_compute_backend_service.api.self_link
  }
}

resource "google_compute_target_https_proxy" "tomedome" {
  name    = "tomedome"
  url_map = google_compute_url_map.tomedome.id
  ssl_certificates = [
    google_compute_managed_ssl_certificate.tomedome_dota.id,
    google_compute_managed_ssl_certificate.tomedome_api.id
  ]
}

resource "google_compute_global_forwarding_rule" "tomedome" {
  name        = "tomedome"
  port_range  = "443"
  ip_protocol = "TCP"
  target      = google_compute_target_https_proxy.tomedome.id
}
