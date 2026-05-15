resource "google_compute_network" "erp_vpc" {
  name                    = "erp-vpc"
  auto_create_subnetworks = false
  project                 = var.project_id
}

resource "google_compute_subnetwork" "erp_subnet" {
  name          = "erp-subnet"
  region        = var.region
  network       = google_compute_network.erp_vpc.id
  ip_cidr_range = "10.0.0.0/24"
  project       = var.project_id
}

resource "google_compute_router" "erp_router" {
  name    = "erp-router"
  region  = var.region
  network = google_compute_network.erp_vpc.id
  project = var.project_id
}

resource "google_compute_router_nat" "erp_nat" {
  name                               = "erp-nat"
  router                             = google_compute_router.erp_router.name
  region                             = var.region
  project                            = var.project_id
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}

resource "google_vpc_access_connector" "erp_connector" {
  name          = "erp-connector"
  region        = var.region
  project       = var.project_id
  network       = google_compute_network.erp_vpc.name
  ip_cidr_range = "10.8.0.0/28"

  min_throughput = 200
  max_throughput = 300
}
