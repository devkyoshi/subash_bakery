output "api_gateway_url" {
  description = "Public HTTPS URL of the API gateway — use as VITE_API_BASE_URL (append /api/v1)"
  value       = module.cloud_run.api_gateway_url
}

output "service_urls" {
  description = "Map of all Cloud Run service names to their HTTPS URIs"
  value       = module.cloud_run.service_urls
}

output "infra_vm_private_ip" {
  description = "Internal VPC IP of the infra VM running MongoDB, Redis, RabbitMQ"
  value       = module.infra_vm.private_ip
}

output "registry_url" {
  description = "Artifact Registry URL prefix for Docker images"
  value       = module.artifact_registry.registry_url
}
