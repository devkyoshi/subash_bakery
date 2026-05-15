locals {
  secrets = {
    "erp-mongo-uri"                  = var.mongo_uri
    "erp-redis-addr"                 = var.redis_addr
    "erp-rabbitmq-url"               = var.rabbitmq_url
    "erp-jwt-secret"                 = var.jwt_secret
    "erp-google-client-id"           = var.google_client_id
    "erp-google-client-secret"       = var.google_client_secret
    "erp-google-redirect-url"        = var.google_redirect_url
    "erp-firebase-credentials-json"  = var.firebase_credentials_json
    "erp-mongo-password"             = var.mongo_password
    "erp-rabbitmq-password"          = var.rabbitmq_password
  }
}

resource "google_secret_manager_secret" "secrets" {
  for_each  = local.secrets
  secret_id = each.key
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "versions" {
  for_each    = local.secrets
  secret      = google_secret_manager_secret.secrets[each.key].id
  secret_data = each.value
}

# Grant Cloud Run service account access to all secrets
resource "google_secret_manager_secret_iam_member" "cloudrun_accessor" {
  for_each  = local.secrets
  project   = var.project_id
  secret_id = google_secret_manager_secret.secrets[each.key].secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${var.cloudrun_sa_email}"
}
