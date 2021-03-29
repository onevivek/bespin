variable "gcp_credentials_file" {
  type        = string
  description = "The GCP credentials file."
}

variable "gcp_project" {
  type        = string
  description = "The GCP project."
}

variable "gcp_region" {
  type        = string
  description = "The GCP region."
}

variable "gcp_zone" {
  type        = string
  description = "The GCP zone."
}

variable "vpc_name" {
  type        = string
  description = "The name of the VPC network."
  default = "devopstest-vpc1"
}
