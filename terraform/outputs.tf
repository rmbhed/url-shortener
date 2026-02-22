output "github_actions_provider_name" {
  value = google_iam_workload_identity_pool_provider.github_provider.name
  description = "The full identifier for the Workload Identity Provider"
}

output "github_actions_service_account_email" {
  value = google_service_account.github_actions_sa.email
}

output "project_number" {
  value = data.google_project.project.number
}

# You'll need this data block at the top of outputs.tf or main.tf
data "google_project" "project" {}