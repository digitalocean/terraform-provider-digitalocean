output "cluster_name" {
  value = module.doks-cluster.primary_cluster.name
}

output "kubeconfig_path" {
  value = var.write_kubeconfig ? abspath("${path.root}/kubeconfig") : "none"
}

output "test_url" {
  value = module.kubernetes-config.test_url
}