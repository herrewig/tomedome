resource "google_artifact_registry_repository" "tomedome" {
  location      = "us-east4"
  repository_id = "tomedome"
  description   = "tomedome.io docker registry"
  format        = "DOCKER"

  cleanup_policy_dry_run = false
  cleanup_policies {
    id     = "keep-minimum-versions"
    action = "KEEP"
    most_recent_versions {
      package_name_prefixes = ["tomedome"]
      keep_count            = 5
    }
  }
}

# This service account will be used to publish images to the registry from CI
resource "google_service_account" "publisher" {
  account_id   = "artifact-registry-writer"
  display_name = "Artifact Registry Writer Service Account"
}

resource "google_project_iam_binding" "publisher" {
  project = data.google_project.current.project_id
  role    = "roles/artifactregistry.writer"

  members = [
    "serviceAccount:${google_service_account.publisher.email}"
  ]
}
