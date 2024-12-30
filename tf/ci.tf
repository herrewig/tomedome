# This service account will be used to publish images to the registry from CI
resource "google_service_account" "ci_admin" {
  account_id   = "ci-admin"
  display_name = "CI Admin"
}

resource "google_project_iam_binding" "publisher" {
  project = data.google_project.current.project_id
  role    = "roles/artifactregistry.writer"

  members = [
    "serviceAccount:${google_service_account.ci_admin.email}"
  ]
}

resource "google_project_iam_binding" "run_admin" {
  project = data.google_project.current.project_id
  role    = "roles/run.admin"

  members = [
    "serviceAccount:${google_service_account.ci_admin.email}"
  ]
}

resource "google_project_iam_binding" "gcs_admin" {
  project = data.google_project.current.project_id
  role    = "roles/storage.admin"

  members = [
    "serviceAccount:${google_service_account.ci_admin.email}"
  ]
}