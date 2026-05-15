# api-gateway is the ONLY service with public (unauthenticated) access
resource "google_cloud_run_service_iam_member" "api_gateway_public" {
  project  = var.project_id
  location = var.region
  service  = google_cloud_run_v2_service.api_gateway.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
