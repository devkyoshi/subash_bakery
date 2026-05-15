output "registry_url" {
  value = "${var.region}-docker.pkg.dev/${var.project_id}/erp-backend"
}

output "repository_name" {
  value = google_artifact_registry_repository.erp_backend.name
}
