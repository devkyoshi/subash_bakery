variable "project_id" {
  type        = string
  description = "GCP project ID (same as Firebase project: subash-bakery)"
}

variable "region" {
  type        = string
  default     = "us-central1"
  description = "GCP region for all resources"
}

variable "image_tag" {
  type        = string
  default     = "latest"
  description = "Docker image tag to deploy. Set to a git SHA for pinned deployments."
}

variable "machine_type" {
  type        = string
  default     = "e2-standard-2"
  description = "GCE machine type for the infra VM. Use e2-medium to reduce cost at low traffic."
}

# Sensitive credentials — set via terraform.tfvars (gitignored) or TF_VAR_* env vars in CI
variable "mongo_password" {
  type      = string
  sensitive = true
}

variable "rabbitmq_password" {
  type      = string
  sensitive = true
}

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

variable "firebase_credentials_json" {
  type        = string
  sensitive   = true
  description = "Full JSON content of firebase-credentials.json"
}

variable "google_redirect_url" {
  type        = string
  default     = "https://placeholder.example.com/api/v1/auth/google/callback"
  description = "Google OAuth redirect URL. Set to https://API_GATEWAY_URL/api/v1/auth/google/callback after first deploy."
}

# Populated after first Cloud Run deploy — needed to break inventory↔product circular dep
variable "product_service_url" {
  type        = string
  default     = ""
  description = "Leave empty on first apply. Set to product-service Cloud Run URL on subsequent applies."
}
