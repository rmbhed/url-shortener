variable "project_id" {
  type        = string
  description = "The Google Cloud Project ID"
}

variable "region" {
  type        = string
  default     = "us-central1"
}

variable "github_repo" {
  type        = string
  description = "The GitHub repository in the format 'owner/repo'"
}