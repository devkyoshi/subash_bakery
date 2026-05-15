# api-gateway accepts public traffic; ingress=ALL is the boundary
resource "google_cloud_run_service_iam_member" "api_gateway_public" {
  project  = var.project_id
  location = var.region
  service  = google_cloud_run_v2_service.api_gateway.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Private services: ingress=INTERNAL_ONLY is the security boundary.
# allUsers invoker lets api-gateway proxy without attaching identity tokens.
locals {
  private_services = {
    auth_service         = google_cloud_run_v2_service.auth_service.name
    org_service          = google_cloud_run_v2_service.org_service.name
    product_service      = google_cloud_run_v2_service.product_service.name
    inventory_service    = google_cloud_run_v2_service.inventory_service.name
    procurement_service  = google_cloud_run_v2_service.procurement_service.name
    notification_service = google_cloud_run_v2_service.notification_service.name
    dashboard_service    = google_cloud_run_v2_service.dashboard_service.name
    report_service       = google_cloud_run_v2_service.report_service.name
  }
}

resource "google_cloud_run_service_iam_member" "private_services_internal_invoker" {
  for_each = local.private_services

  project  = var.project_id
  location = var.region
  service  = each.value
  role     = "roles/run.invoker"
  member   = "allUsers"
}
