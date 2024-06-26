provider "awsssmtunnels" {
  region     = "us-east-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
  target     = "i-123456789"
}

// OR

provider "awsssmtunnels" {
  region     = "us-east-1"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
  token      = var.aws_token
  target     = "i-123456789"
}

// OR
provider "awsssmtunnels" {
  region              = "us-east-1"
  shared_config_files = [var.tfc_aws_dynamic_credentials.default.shared_config_file]
  target              = "i-123456789"
}
