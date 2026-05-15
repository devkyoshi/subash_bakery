variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "cicd_sa_email" {
  type        = string
  description = "Email of the GitHub Actions CI/CD service account"
}

variable "cloudrun_sa_email" {
  type        = string
  description = "Email of the Cloud Run service account"
}
