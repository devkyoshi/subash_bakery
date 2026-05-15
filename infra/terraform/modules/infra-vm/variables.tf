variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "us-central1"
}

variable "subnet_id" {
  type        = string
  description = "Subnetwork ID for the VM network interface"
}

variable "machine_type" {
  type        = string
  default     = "e2-standard-2"
  description = "GCE machine type. Use e2-medium for lower traffic."
}

variable "mongo_password" {
  type      = string
  sensitive = true
}

variable "rabbitmq_password" {
  type      = string
  sensitive = true
}

variable "vm_sa_email" {
  type        = string
  description = "Service account email to attach to the VM"
}
