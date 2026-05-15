resource "google_compute_instance" "erp_infra" {
  name         = "erp-infra"
  machine_type = var.machine_type
  zone         = "${var.region}-a"
  project      = var.project_id

  tags = ["erp-infra"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
      size  = 50
      type  = "pd-ssd"
    }
  }

  network_interface {
    subnetwork = var.subnet_id
    # No access_config = no public IP (IAP SSH only)
  }

  metadata = {
    startup-script         = templatefile("${path.module}/startup.sh", {
      mongo_password    = var.mongo_password
      rabbitmq_password = var.rabbitmq_password
    })
    enable-oslogin         = "TRUE"
    serial-port-enable     = "FALSE"
  }

  service_account {
    email  = var.vm_sa_email
    scopes = ["cloud-platform"]
  }

  allow_stopping_for_update = true
}

# Service account for the infra VM (needs Secret Manager read access)
resource "google_service_account" "erp_vm_sa" {
  account_id   = "erp-infra-vm-sa"
  display_name = "ERP Infra VM Service Account"
  project      = var.project_id
}

resource "google_project_iam_member" "vm_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.erp_vm_sa.email}"
}
