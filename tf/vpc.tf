resource "google_compute_network" "tomedome" {
  name                    = "tomedome"
  auto_create_subnetworks = false
}

# Small and simple
resource "google_compute_subnetwork" "tomedome" {
  name          = "tomedome"
  network       = google_compute_network.tomedome.self_link
  ip_cidr_range = "172.25.100.0/24"
}