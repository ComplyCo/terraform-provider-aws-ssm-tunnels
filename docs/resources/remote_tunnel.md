---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awsssmtunnels_remote_tunnel Resource - awsssmtunnels"
subcategory: ""
description: |-
  AWSM SSM Remote Tunnel data source
---

# awsssmtunnels_remote_tunnel (Resource)

AWSM SSM Remote Tunnel data source

## Example Usage

```terraform
##############################################
######## Prerequisites #######################
##############################################

// Make sure you have one or more bastion instances running in the correct VPCs and with the correct
// Security groups to allow connectivity to EKS/RDS and whatever else you are trying to connect to.
// In the examples below, we have a single bastion instance with the ID i-123456789. The instance
// has 2 security groups applied. One allows port 443 outbound to EKS (and EKS's security group allows
// 443 inbound from the bastion's security group). The bastion has another security group that allows port
// 5432 outbound to RDS (and RDS's security group allows 5432 inbound from the bastion's security group).

##############################################
######## EKS Example #########################
##############################################

provider "kubernetes" {
  host                   = "https://${awsssmtunnels_remote_tunnel.eks.local_host}:${awsssmtunnels_remote_tunnel.eks.local_port}"
  tls_server_name        = replace(aws_eks_cluster.example.endpoint, "https://", "")
  cluster_ca_certificate = base64decode(aws_eks_cluster.example.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.example.token
}

data "aws_eks_cluster_auth" "example" {
  name = aws_eks_cluster.example.name
}

resource "awsssmtunnels_remote_tunnel" "eks" {
  refresh_id  = "one" // Anything string can go here as this resource will always find a diff on this
  target      = "i-123456789"
  remote_host = replace(aws_eks_cluster.example.endpoint, "https://", "")
  remote_port = 443
  local_port  = 16534
  region      = "us-east-1"
}

// NOTE: The import is needed for the first plan, otherwise TF will hold off on the create until the Apply phase.
// The import allows it to run in the very first plan.
import {
  id = "i-123456789|${replace(aws_eks_cluster.example.endpoint, "https://", "")}|443|16534|127.0.0.1|us-east-1"
  to = awsssmtunnels_remote_tunnel.eks
}

// NOTE: We use the *_keepalive data resource to prevent the provider from being shut down prematurely
// We need the tunnel to stay up until all the resources for the providers using the tunnel are done
// reading or writing from it.
data "awsssmtunnels_keepalive" "eks" {
  depends_on = [
    kubernetes_secret.one,
    kubernetes_secret.two,
    kubernetes_config_map.one,
    kubernetes_config_map.two,
    helm_release.example_operator,
    awsssmtunnels_remote_tunnel.eks,
  ]
}


##############################################
######## RDS Example #########################
##############################################

resource "awsssmtunnels_remote_tunnel" "rds" {
  refresh_id  = "one" // Anything string can go here as this resource will always find a diff on this
  remote_host = aws_rds_cluster.example.endpoint
  remote_port = 5432 // This is a PostgreSQL RDS cluster example
  local_port  = 17638
}

import {
  id = "<remote host>|<remote port>|<local port>|127.0.0.1"
  to = awsssmtunnels_remote_tunnel.rds
}

data "awsssmtunnels_keepalive" "rds" {
  depends_on = [
    postgresql_tables.my_tables,
    awsssmtunnels_remote_tunnel.rds,
  ]
}

provider "postgresql" {
  host            = awsssmtunnels_remote_tunnel.rds.local_host
  port            = awsssmtunnels_remote_tunnel.rds.local_port
  database        = "mydb"
  username        = var.pg_user
  password        = var.pg_password
  sslmode         = "require"
  connect_timeout = 15
}

data "postgresql_tables" "my_tables" {
  database   = "mydb"
  depends_on = [awsssmtunnels_remote_tunnel.rds] // NOTE: The tunnel must be up before we can query the database
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `refresh_id` (String) Any value as this will trigger a refresh
- `remote_host` (String) The DNS name or IP address of the remote host
- `remote_port` (Number) The port number of the remote host

### Optional

- `local_port` (Number) The local port number to use for the tunnel

### Read-Only

- `id` (String) Example identifier
- `local_host` (String) The DNS name or IP address of the local host
