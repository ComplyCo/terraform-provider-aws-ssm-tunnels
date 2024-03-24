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
  ]
}

data "awsssmtunnels_keepalive" "rds" {
  depends_on = [
    postgresql_tables.my_tables,
  ]
}
