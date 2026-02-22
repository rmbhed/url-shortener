terraform {
  backend "gcs" {
    bucket = "rmbh-url-shortener-tfstate"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Create the Artifact Registry repository
resource "google_artifact_registry_repository" "shortener_repo" {
  location      = var.region
  repository_id = "shortener-app"
  description   = "Docker repository for the URL shortener"
  format        = "DOCKER"
}

# 1. Enable the necessary Google APIs
resource "google_project_service" "run_api" {
  service = "run.googleapis.com"
}

resource "google_project_service" "firestore_api" {
  service = "firestore.googleapis.com"
}

# 2. Create the Firestore Database
resource "google_firestore_database" "database" {
  project     = var.project_id
  name        = "(default)"
  location_id = var.region
  type        = "FIRESTORE_NATIVE"

  depends_on = [google_project_service.firestore_api]
}

# 3. Cloud Run Service
resource "google_cloud_run_v2_service" "url_shortener" {
  name     = "url-shortener-service"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello" 
      
      env {
        name  = "PROJECT_ID"
        value = var.project_id
      }
    }
  }

  depends_on = [google_project_service.run_api]
}

# 4. Domain Mapping (The "short.rmbh.me" magic)
resource "google_cloud_run_domain_mapping" "custom_domain" {
  location = var.region
  name     = "short.rmbh.me"

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.url_shortener.name
  }
}

# 5. Make it Public
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  name     = google_cloud_run_v2_service.url_shortener.name
  location = google_cloud_run_v2_service.url_shortener.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}