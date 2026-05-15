# Allow MongoDB, Redis, RabbitMQ only from within the VPC subnet (Cloud Run via connector)
resource "google_compute_firewall" "allow_infra_internal" {
  name     = "allow-infra-internal"
  network  = google_compute_network.erp_vpc.name
  project  = var.project_id
  priority = 900

  allow {
    protocol = "tcp"
    ports    = ["27017", "6379", "5672", "15672"]
  }

  source_ranges = ["10.0.0.0/24", "10.8.0.0/28"]
  target_tags   = ["erp-infra"]
}

# Allow SSH to the infra VM from IAP (Identity-Aware Proxy) only — no public SSH
resource "google_compute_firewall" "allow_iap_ssh" {
  name    = "allow-iap-ssh"
  network = google_compute_network.erp_vpc.name
  project = var.project_id

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["35.235.240.0/20"]
  target_tags   = ["erp-infra"]
}

# Deny all other external ingress to infra VM
resource "google_compute_firewall" "deny_external_infra" {
  name     = "deny-external-infra"
  network  = google_compute_network.erp_vpc.name
  project  = var.project_id
  priority = 1000

  deny {
    protocol = "tcp"
    ports    = ["27017", "6379", "5672"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["erp-infra"]
}

# Allow egress from VPC (so VM can pull Docker images)
resource "google_compute_firewall" "allow_egress" {
  name      = "allow-egress"
  network   = google_compute_network.erp_vpc.name
  project   = var.project_id
  direction = "EGRESS"

  allow {
    protocol = "all"
  }

  destination_ranges = ["0.0.0.0/0"]
}
