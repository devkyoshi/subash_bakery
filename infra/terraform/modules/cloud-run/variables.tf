variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "image_tag" {
  type        = string
  default     = "latest"
  description = "Docker image tag to deploy across all services"
}

variable "connector_id" {
  type        = string
  description = "Serverless VPC Connector ID for private network access"
}

variable "cloudrun_sa_email" {
  type        = string
  description = "Email of the Cloud Run service account (created in root module)"
}

variable "secret_ids" {
  type        = map(string)
  description = "Map of secret name → Secret Manager secret_id"
}

# Circular dependency workaround: inventory needs product URL but product needs inventory URL.
# Pass product-service URL as a variable (set after first deployment, then re-apply).
variable "product_service_url" {
  type        = string
  default     = ""
  description = "Cloud Run URI for product-service. Set after first deploy to break circular dependency."
}
