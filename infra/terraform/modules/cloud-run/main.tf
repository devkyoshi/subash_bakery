# Helper local: build a secret env var source block
locals {
  registry = "${var.region}-docker.pkg.dev/${var.project_id}/erp-backend"

  # Common env vars shared by all services
  common_env = [
    { name = "ENV", value = "production" },
    { name = "LOG_LEVEL", value = "info" },
  ]

  # Secret-backed env vars shared by all services
  common_secrets = {
    "MONGO_URI"    = var.secret_ids["erp-mongo-uri"]
    "REDIS_ADDR"   = var.secret_ids["erp-redis-addr"]
    "RABBITMQ_URL" = var.secret_ids["erp-rabbitmq-url"]
    "JWT_SECRET"   = var.secret_ids["erp-jwt-secret"]
  }
}

# ── auth-service ──────────────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "auth_service" {
  name     = "auth-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = "${local.registry}/auth-service:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      env {
        name = "GOOGLE_CLIENT_ID"
        value_source {
          secret_key_ref {
            secret  = var.secret_ids["erp-google-client-id"]
            version = "latest"
          }
        }
      }

      env {
        name = "GOOGLE_CLIENT_SECRET"
        value_source {
          secret_key_ref {
            secret  = var.secret_ids["erp-google-client-secret"]
            version = "latest"
          }
        }
      }

      env {
        name = "GOOGLE_REDIRECT_URL"
        value_source {
          secret_key_ref {
            secret  = var.secret_ids["erp-google-redirect-url"]
            version = "latest"
          }
        }
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── org-service ───────────────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "org_service" {
  name     = "org-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = "${local.registry}/org-service:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── product-service ───────────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "product_service" {
  name     = "product-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = "${local.registry}/product-service:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      env {
        name  = "INVENTORY_SERVICE_URL"
        value = google_cloud_run_v2_service.inventory_service.uri
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  depends_on = [google_cloud_run_v2_service.inventory_service]

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── inventory-service ─────────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "inventory_service" {
  name     = "inventory-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = "${local.registry}/inventory-service:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      env {
        name  = "AUTH_SERVICE_URL"
        value = google_cloud_run_v2_service.auth_service.uri
      }

      env {
        name  = "ORG_SERVICE_URL"
        value = google_cloud_run_v2_service.org_service.uri
      }

      # product-service URL resolved at runtime via env var update after product-service deploys
      env {
        name  = "PRODUCT_SERVICE_URL"
        value = var.product_service_url
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── procurement-service ───────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "procurement_service" {
  name     = "procurement-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = "${local.registry}/procurement-service:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      env {
        name  = "AUTH_SERVICE_URL"
        value = google_cloud_run_v2_service.auth_service.uri
      }

      env {
        name  = "PRODUCT_SERVICE_URL"
        value = google_cloud_run_v2_service.product_service.uri
      }

      env {
        name  = "INVENTORY_SERVICE_URL"
        value = google_cloud_run_v2_service.inventory_service.uri
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  depends_on = [
    google_cloud_run_v2_service.auth_service,
    google_cloud_run_v2_service.product_service,
    google_cloud_run_v2_service.inventory_service,
  ]

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── notification-service ──────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "notification_service" {
  name     = "notification-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 3
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    # Mount firebase credentials JSON from Secret Manager as a file
    volumes {
      name = "firebase-credentials"
      secret {
        secret = var.secret_ids["erp-firebase-credentials-json"]
        items {
          version = "latest"
          path    = "firebase-credentials.json"
          mode    = 0444
        }
      }
    }

    containers {
      image = "${local.registry}/notification-service:${var.image_tag}"

      volume_mounts {
        name       = "firebase-credentials"
        mount_path = "/secrets"
      }

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      env {
        name  = "FIREBASE_CREDENTIALS_PATH"
        value = "/secrets/firebase-credentials.json"
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── dashboard-service ─────────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "dashboard_service" {
  name     = "dashboard-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = "${local.registry}/dashboard-service:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── report-service ────────────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "report_service" {
  name     = "report-service"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    vpc_access {
      connector = var.connector_id
      egress    = "PRIVATE_RANGES_ONLY"
    }

    containers {
      image = "${local.registry}/report-service:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      dynamic "env" {
        for_each = local.common_secrets
        content {
          name = env.key
          value_source {
            secret_key_ref {
              secret  = env.value
              version = "latest"
            }
          }
        }
      }

      env {
        name  = "PROCUREMENT_SERVICE_URL"
        value = google_cloud_run_v2_service.procurement_service.uri
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  depends_on = [google_cloud_run_v2_service.procurement_service]

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}

# ── api-gateway ───────────────────────────────────────────────────────────────
resource "google_cloud_run_v2_service" "api_gateway" {
  name     = "api-gateway"
  location = var.region
  project  = var.project_id
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = var.cloudrun_sa_email

    scaling {
      min_instance_count = 1
      max_instance_count = 10
    }

    vpc_access {
      connector = var.connector_id
      egress    = "ALL_TRAFFIC"
    }

    containers {
      image = "${local.registry}/api-gateway:${var.image_tag}"

      dynamic "env" {
        for_each = local.common_env
        content {
          name  = env.value.name
          value = env.value.value
        }
      }

      env {
        name  = "AUTH_SERVICE_URL"
        value = google_cloud_run_v2_service.auth_service.uri
      }

      env {
        name  = "ORG_SERVICE_URL"
        value = google_cloud_run_v2_service.org_service.uri
      }

      env {
        name  = "PRODUCT_SERVICE_URL"
        value = google_cloud_run_v2_service.product_service.uri
      }

      env {
        name  = "INVENTORY_SERVICE_URL"
        value = google_cloud_run_v2_service.inventory_service.uri
      }

      env {
        name  = "PROCUREMENT_SERVICE_URL"
        value = google_cloud_run_v2_service.procurement_service.uri
      }

      env {
        name  = "NOTIFICATION_SERVICE_URL"
        value = google_cloud_run_v2_service.notification_service.uri
      }

      env {
        name  = "DASHBOARD_SERVICE_URL"
        value = google_cloud_run_v2_service.dashboard_service.uri
      }

      env {
        name  = "REPORT_SERVICE_URL"
        value = google_cloud_run_v2_service.report_service.uri
      }

      resources {
        limits = {
          cpu    = "2"
          memory = "1Gi"
        }
      }
    }
  }

  depends_on = [
    google_cloud_run_v2_service.auth_service,
    google_cloud_run_v2_service.org_service,
    google_cloud_run_v2_service.product_service,
    google_cloud_run_v2_service.inventory_service,
    google_cloud_run_v2_service.procurement_service,
    google_cloud_run_v2_service.notification_service,
    google_cloud_run_v2_service.dashboard_service,
    google_cloud_run_v2_service.report_service,
  ]

  lifecycle {
    ignore_changes = [template[0].containers[0].image]
  }
}
