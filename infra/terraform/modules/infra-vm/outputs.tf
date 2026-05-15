output "private_ip" {
  value       = google_compute_instance.erp_infra.network_interface[0].network_ip
  description = "Internal VPC IP of the infra VM — used to build connection strings for secrets"
}

output "vm_sa_email" {
  value = google_service_account.erp_vm_sa.email
}
