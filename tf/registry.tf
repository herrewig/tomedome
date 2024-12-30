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
