output "secret_ids" {
  value = {
    for k, v in google_secret_manager_secret.secrets : k => v.secret_id
  }
}

output "mongo_uri_secret_id" {
  value = google_secret_manager_secret.secrets["erp-mongo-uri"].secret_id
}

output "redis_addr_secret_id" {
  value = google_secret_manager_secret.secrets["erp-redis-addr"].secret_id
}

output "rabbitmq_url_secret_id" {
  value = google_secret_manager_secret.secrets["erp-rabbitmq-url"].secret_id
}

output "jwt_secret_id" {
  value = google_secret_manager_secret.secrets["erp-jwt-secret"].secret_id
}

output "google_client_id_secret_id" {
  value = google_secret_manager_secret.secrets["erp-google-client-id"].secret_id
}

output "google_client_secret_secret_id" {
  value = google_secret_manager_secret.secrets["erp-google-client-secret"].secret_id
}

output "google_redirect_url_secret_id" {
  value = google_secret_manager_secret.secrets["erp-google-redirect-url"].secret_id
}

output "firebase_credentials_secret_id" {
  value = google_secret_manager_secret.secrets["erp-firebase-credentials-json"].secret_id
}
