# #### #
# Data Sources
# #### #
data "google_client_config" "default" {}

# #### #
# Variables
# #### #
# The following four variables will be set by environment variables
variable "gcp_credentials_file" {}
variable "gcp_project" {}
variable "gcp_region" {}
variable "gcp_zone" {}
variable "cluster_name" { default = "devopstest-gke1" }
variable "network" { default = "default" }
variable "subnetwork" { default = "" }
variable "ip_range_pods" { default = "" }
variable "ip_range_services" { default = "" }

# #### #
# Providers
# #### #
# Set GCP settings using values from environment variables
provider "google" {
  credentials = file(var.gcp_credentials_file)
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}
# Set Kubernetes provider
provider "kubernetes" {
  load_config_file       = false
  host                   = "https://${module.gke.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(module.gke.ca_certificate)
}

# #### #
# GKE Cluster
#   - No VPC, No autoscaling, zonal (not regional)
# #### #
module "gke" {
  # #### #
  # Source and source version
  # #### #
  source  = "terraform-google-modules/kubernetes-engine/google"
  version = "14.0.1"

  # #### #
  # Required variables
  # #### #
  project_id        = var.gcp_project
  name              = var.cluster_name
  region            = var.gcp_region
  zones             = [var.gcp_region]
  network           = var.network
  subnetwork        = var.subnetwork
  ip_range_pods     = var.ip_range_pods
  ip_range_services = var.ip_range_services

  # #### #
  # Optional variables
  # #### #
  kubernetes_version       = "1.18.16-gke.1200"
  regional                 = false
  create_service_account   = false
  remove_default_node_pool = true

  # #### #
  # Addons
  # #### #
  network_policy             = false
  horizontal_pod_autoscaling = false
  http_load_balancing        = false

  # #### #
  # Node Pools
  # #### #
  node_pools = [
    {
      name               = "default-node-pool"
      min_count          = 2
      max_count          = 4
      initial_node_count = 1
      node_locations     = var.gcp_zone
      machine_type       = "n1-standard-1"
      local_ssd_count    = 0
      disk_size_gb       = 100
      disk_type          = "pd-standard"
      image_type         = "COS"
      auto_repair        = true
      auto_upgrade       = true
      preemptible        = false
    },
  ]

  node_pools_oauth_scopes = {
    all = []
    default-node-pool = [
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/trace.append",
    ]
  }

  node_pools_labels = {
    all = {}
    default-node-pool = {
      default-node-pool = true,
    }
  }

  node_pools_tags = {
    all = []
    default-node-pool = [
      "default-node-pool",
    ]
  }
}
