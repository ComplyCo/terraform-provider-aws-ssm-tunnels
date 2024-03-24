terraform {
  required_version = "~> 1.7"

  required_providers {
    cc = {
      source  = "app.terraform.io/ComplyCo/aws-ssm-tunnels"
      version = "0.0.8"
    }
  }
}

provider "cc" {
  region     = "us-east-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}


provider "cc" {
  region     = "us-east-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
  token      = var.aws_token
}

provider "cc" {
  region              = "us-east-1"
  shared_config_files = [var.tfc_aws_dynamic_credentials.default.shared_config_file]
}
