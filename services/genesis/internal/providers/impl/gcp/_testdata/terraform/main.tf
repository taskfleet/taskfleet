terraform {
  backend "local" {}
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.42"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.4"
    }
  }
}

variable "gcp_project" {}

provider "google" {
  project = var.gcp_project
}

#--------------------------------------------------------------------------------------------------

variable "create_iam" {
  default = false
}

#--------------------------------------------------------------------------------------------------

output "service_account_email" {
  value = var.create_iam ? google_service_account.test[0].email : ""
}

output "network_name" {
  value = google_compute_network.test.name
}

#--------------------------------------------------------------------------------------------------
# RANDOM VALUE FOR CONCURRENT TESTS
#--------------------------------------------------------------------------------------------------

resource "random_uuid" "test" {}

#--------------------------------------------------------------------------------------------------
# NETWORK
#--------------------------------------------------------------------------------------------------

resource "google_compute_network" "test" {
  name                    = "testnet-${random_uuid.test.result}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "europe_west3" {
  network       = google_compute_network.test.id
  name          = "subnet-${random_uuid.test.result}"
  region        = "europe-west3"
  ip_cidr_range = "10.0.0.0/24"
}

resource "google_compute_subnetwork" "europe_north1" {
  network       = google_compute_network.test.id
  name          = "subnet-${random_uuid.test.result}"
  region        = "us-east1"
  ip_cidr_range = "10.0.1.0/24"
}

#--------------------------------------------------------------------------------------------------
# IAM
#--------------------------------------------------------------------------------------------------

resource "google_service_account" "test" {
  count = var.create_iam ? 1 : 0
  # Need to substring the UUID as it is too long
  account_id = "test-${substr(random_uuid.test.result, 0, 23)}"
}
