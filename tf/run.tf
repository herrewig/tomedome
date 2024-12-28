resource "google_cloud_run_v2_service" "api" {
  name                = "tomedome-api"
  location            = "us-east4"
  deletion_protection = false
  ingress             = "INGRESS_TRAFFIC_ALL"

  template {
    max_instance_request_concurrency = "120"
    containers {
      ports {
        container_port = 8080
      }
      image = "us-east4-docker.pkg.dev/tomedome/tomedome/api:latest"
      liveness_probe {
        http_get {
          path = "/api/v1/healthz"
        }
      }
      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
    scaling {
      max_instance_count = 2
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "api" {
  location = "us-east4"
  name     = google_cloud_run_v2_service.api.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_compute_region_network_endpoint_group" "api" {
  name                  = "tomedome-api"
  network_endpoint_type = "SERVERLESS"
  region                = "us-east4"

  cloud_run {
    service = google_cloud_run_v2_service.api.name
  }
}

resource "google_compute_backend_service" "api" {
  name                  = "tomedome-api"
  protocol              = "HTTP"
  load_balancing_scheme = "EXTERNAL"

  backend {
    group = google_compute_region_network_endpoint_group.api.id
  }
}

resource "google_compute_health_check" "healthz" {
  name                = "healthz"
  check_interval_sec  = 5
  timeout_sec         = 5
  healthy_threshold   = 1
  unhealthy_threshold = 2

  http_health_check {
    port         = 8080
    request_path = "/healthz"
  }
}
