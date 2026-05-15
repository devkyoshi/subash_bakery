output "api_gateway_url" {
  value = google_cloud_run_v2_service.api_gateway.uri
}

output "service_urls" {
  value = {
    "api-gateway"          = google_cloud_run_v2_service.api_gateway.uri
    "auth-service"         = google_cloud_run_v2_service.auth_service.uri
    "org-service"          = google_cloud_run_v2_service.org_service.uri
    "product-service"      = google_cloud_run_v2_service.product_service.uri
    "inventory-service"    = google_cloud_run_v2_service.inventory_service.uri
    "procurement-service"  = google_cloud_run_v2_service.procurement_service.uri
    "notification-service" = google_cloud_run_v2_service.notification_service.uri
    "dashboard-service"    = google_cloud_run_v2_service.dashboard_service.uri
    "report-service"       = google_cloud_run_v2_service.report_service.uri
  }
}

output "cloudrun_sa_email" {
  value = var.cloudrun_sa_email
}
