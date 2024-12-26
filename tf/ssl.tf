resource "google_compute_managed_ssl_certificate" "tomedome_dota" {
  name = "tomedome-dota"

  managed {
    domains = ["dota.tomedome.io"]
  }
}

resource "google_compute_managed_ssl_certificate" "tomedome_api" {
  name = "tomedome-api"

  managed {
    domains = ["api.tomedome.io"]
  }
}
