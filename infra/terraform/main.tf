# Create the Cloud Run service account at root level so both
# the secrets module and cloud-run module can reference it without a cycle.
resource "google_service_account" "cloudrun_sa" {
  account_id   = "erp-cloud-run-sa"
  display_name = "ERP Cloud Run Service Account"
  project      = var.project_id
}

resource "google_project_iam_member" "cloudrun_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.cloudrun_sa.email}"
}

module "networking" {
  source     = "./modules/networking"
  project_id = var.project_id
  region     = var.region
}

module "infra_vm" {
  source     = "./modules/infra-vm"
  project_id = var.project_id
  region     = var.region
  subnet_id  = module.networking.subnet_id

  machine_type      = var.machine_type
  mongo_password    = var.mongo_password
  rabbitmq_password = var.rabbitmq_password
  vm_sa_email       = module.infra_vm.vm_sa_email

  depends_on = [module.networking]
}

module "artifact_registry" {
  source     = "./modules/artifact-registry"
  project_id = var.project_id
  region     = var.region

  cicd_sa_email     = "github-actions@${var.project_id}.iam.gserviceaccount.com"
  cloudrun_sa_email = google_service_account.cloudrun_sa.email
}

module "secrets" {
  source     = "./modules/secrets"
  project_id = var.project_id

  cloudrun_sa_email = google_service_account.cloudrun_sa.email

  mongo_uri    = "mongodb://admin:${var.mongo_password}@${module.infra_vm.private_ip}:27017/erp_db?authSource=admin"
  redis_addr   = "${module.infra_vm.private_ip}:6379"
  rabbitmq_url = "amqp://admin:${var.rabbitmq_password}@${module.infra_vm.private_ip}:5672/"

  jwt_secret                = var.jwt_secret
  google_client_id          = var.google_client_id
  google_client_secret      = var.google_client_secret
  google_redirect_url       = var.google_redirect_url
  firebase_credentials_json = var.firebase_credentials_json
  mongo_password            = var.mongo_password
  rabbitmq_password         = var.rabbitmq_password

  depends_on = [module.infra_vm]
}

module "cloud_run" {
  source     = "./modules/cloud-run"
  project_id = var.project_id
  region     = var.region
  image_tag  = var.image_tag

  cloudrun_sa_email   = google_service_account.cloudrun_sa.email
  connector_id        = module.networking.connector_id
  secret_ids          = module.secrets.secret_ids
  product_service_url = var.product_service_url

  depends_on = [module.networking, module.secrets]
}
