variable "services" {
  type = list(string)
  default = [
    "artifactregistry.googleapis.com",
    "storage-component.googleapis.com",
    "run.googleapis.com",
    "compute.googleapis.com",
    "dns.googleapis.com",
  ]
}

resource "google_project_service" "target" {
  for_each = toset(var.services)

  service = each.value
  disable_on_destroy = true
}