output "vpc_id" {
  value = google_compute_network.erp_vpc.id
}

output "vpc_name" {
  value = google_compute_network.erp_vpc.name
}

output "subnet_id" {
  value = google_compute_subnetwork.erp_subnet.id
}

output "subnet_name" {
  value = google_compute_subnetwork.erp_subnet.name
}

output "connector_id" {
  value = google_vpc_access_connector.erp_connector.id
}
