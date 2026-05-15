variable "project_id" {
  type = string
}

variable "cloudrun_sa_email" {
  type        = string
  description = "Email of the Cloud Run service account that needs secret access"
}

# Infrastructure connection strings (built from infra-vm outputs)
variable "mongo_uri" {
  type      = string
  sensitive = true
}

variable "redis_addr" {
  type      = string
  sensitive = true
}

variable "rabbitmq_url" {
  type      = string
  sensitive = true
}

# Application secrets
variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "google_client_id" {
  type      = string
  sensitive = true
}

variable "google_client_secret" {
  type      = string
  sensitive = true
}

variable "google_redirect_url" {
  type = string
}

variable "firebase_credentials_json" {
  type        = string
  sensitive   = true
  description = "Full JSON content of firebase-credentials.json"
}

# Raw passwords used by startup script
variable "mongo_password" {
  type      = string
  sensitive = true
}

variable "rabbitmq_password" {
  type      = string
  sensitive = true
}
