terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {
  credentials = file(var.gcp_credentials_file)
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

resource "google_compute_network" "vpc_network" {
  name = var.vpc_name
}
