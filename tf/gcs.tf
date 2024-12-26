resource "google_storage_bucket" "static_site" {
  name          = "tomedome-static-site"
  location      = "us-east4"
  storage_class = "REGIONAL"
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
  }
}

resource "google_storage_bucket_iam_binding" "public_access" {
  bucket = google_storage_bucket.static_site.name

  role = "roles/storage.objectViewer"
  members = [
    "allUsers",
  ]
}

resource "google_compute_backend_bucket" "static_site_backend" {
  name        = "tomedome-static-site-backend"
  bucket_name = google_storage_bucket.static_site.name
  enable_cdn  = false
}

