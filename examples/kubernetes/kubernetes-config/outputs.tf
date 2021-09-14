output "test_url" {
  value = "http://${data.kubernetes_service.nginx-ingress-controller.status.0.load_balancer.0.ingress.0.ip}/test"
}
