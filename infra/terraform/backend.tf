terraform {
  backend "gcs" {
    bucket = "subash-bakery-495502-tfstate"
    prefix = "terraform/state"
  }
}
