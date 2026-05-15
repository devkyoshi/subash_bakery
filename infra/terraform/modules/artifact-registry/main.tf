resource "google_artifact_registry_repository" "erp_backend" {
  location      = var.region
  repository_id = "erp-backend"
  format        = "DOCKER"
  project       = var.project_id
  description   = "Docker images for ERP backend microservices"
}

# Allow CI/CD service account to push images
resource "google_artifact_registry_repository_iam_member" "cicd_writer" {
  project    = var.project_id
  location   = var.region
  repository = google_artifact_registry_repository.erp_backend.name
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${var.cicd_sa_email}"
}

# Allow Cloud Run service account to pull images
resource "google_artifact_registry_repository_iam_member" "cloudrun_reader" {
  project    = var.project_id
  location   = var.region
  repository = google_artifact_registry_repository.erp_backend.name
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${var.cloudrun_sa_email}"
}
