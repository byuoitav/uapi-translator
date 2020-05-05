terraform {
  backend "s3" {
    bucket     = "terraform-state-storage-586877430255"
    lock_table = "terraform-state-lock-586877430255"
    region     = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "av-uapi.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
}

// pull all env vars out of ssm

data "aws_ssm_parameter" "auth_url" {
  name = "/env/av-uapi/opa-url"
}

data "aws_ssm_parameter" "uapi_auth_token" {
  name = "/env/av-uapi/auth-token"
}

data "aws_ssm_parameter" "prd_av_api_url" {
  name = "/env/av-uapi/av-api-url"
}

data "aws_ssm_parameter" "prd_db_address" {
  name = "/env/couch-address"
}

data "aws_ssm_parameter" "prd_db_password" {
  name = "/env/couch-password"
}

data "aws_ssm_parameter" "prd_db_username" {
  name = "/env/couch-username"
}

module "av_uapi_prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "av-uapi-prd"
  image          = "byuoitav/av-uapi"
  image_version  = "production"
  container_port = 80
  repo_url       = "https://github.com/byuoitav/uapi-translator"

  // optional
  public_urls = ["uapi.av.byu.edu"]
  container_env = {
    DB_ADDRESS  = data.aws_ssm_parameter.prd_db_address.value
    DB_PASSWORD = data.aws_ssm_parameter.prd_db_password.value
    DB_USERNAME = data.aws_ssm_parameter.prd_db_username.value
    AV_API_URL  = data.aws_ssm_parameter.prd_av_api_url.value
  }
  container_args = [
    "--opa-url", data.aws_ssm_parameter.auth_url.value,
    "--opa-token", data.aws_ssm_parameter.uapi_auth_token.value,
    "--db-address", data.aws_ssm_parameter.prd_db_address.value,
    "--db-username", data.aws_ssm_parameter.prd_db_username.value,
    "--db-password", data.aws_ssm_parameter.prd_db_password.value,
  ]
}
