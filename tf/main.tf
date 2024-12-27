terraform {
  required_version = ">= 1.10.1"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.13.0"
    }
  }

  backend "gcs" {
    # This bucket needs to be created manually for bootstrapping
    bucket = "tomedome-io-tf-state"
    prefix = "terraform/state"
  }

}

provider "google" {
  project = "tomedome"
  region  = "us-east4"
}

module "gcp_services" {
  source = "./modules/services"
}

data "google_project" "current" {}
