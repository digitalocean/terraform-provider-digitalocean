terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.4.0"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.0.1"
    }
  }
}

resource "local_file" "kubeconfig" {
  depends_on = [var.primary_cluster]
  count      = var.write_kubeconfig ? 1 : 0
  content    = var.primary_cluster.raw_config
  filename   = "${path.root}/kubeconfig"
}

provider "kubernetes" {
  host                   = var.primary_cluster.endpoint
  token                  = var.primary_cluster.token
  cluster_ca_certificate = base64decode(var.primary_cluster.cluster_ca_certificate)
}

provider "helm" {
  kubernetes {
    host                   = var.primary_cluster.endpoint
    token                  = var.primary_cluster.token
    cluster_ca_certificate = base64decode(var.primary_cluster.cluster_ca_certificate)
  }
}

resource "kubernetes_namespace" "test" {
  metadata {
    name = "test"
  }
}

resource "kubernetes_deployment" "test" {
  metadata {
    name = "test"
    namespace= kubernetes_namespace.test.metadata.0.name
  }
  spec {
    replicas = 2
    selector {
      match_labels = {
        app = "test"
      }
    }
    template {
      metadata {
        labels = {
          app  = "test"
        }
      }
      spec {
        container {
          image = "hashicorp/http-echo"
          name  = "http-echo"
          args  = ["-text=test"]

          resources {
            limits = {
              memory = "512M"
              cpu = "1"
            }
            requests = {
              memory = "256M"
              cpu = "50m"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "test" {
  metadata {
    name      = "test-service"
    namespace = kubernetes_namespace.test.metadata.0.name
  }
  spec {
    selector = {
      app = kubernetes_deployment.test.metadata.0.name
    }

    port {
      port = 5678
    }
  }
}

resource "helm_release" "nginx_ingress" {
  name       = "nginx-ingress-controller"
  namespace  = kubernetes_namespace.test.metadata.0.name

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx-ingress-controller"

  # Helm chart deployment can sometimes take longer than the default 5 minutes
  timeout    = var.nginx_ingress_helm_timeout_seconds

  # Try to allow time for external endpoints to be applied to service
  wait_for_jobs = true

  set {
    name  = "service.type"
    value = "LoadBalancer"
  }
  set {
    name  = "service.annotations.service\\.beta\\.kubernetes\\.io/do-loadbalancer-name"
    value = var.primary_cluster.name
  }
  set {
    name  = "service.annotations.service\\.beta\\.kubernetes\\.io/do-loadbalancer-size-slug"
    value = "lb-small"
  }
}

resource "kubernetes_ingress" "test_ingress" {
  wait_for_load_balancer = true
  metadata {
    name = "test-ingress"
    namespace  = kubernetes_namespace.test.metadata.0.name
    annotations = {
      "kubernetes.io/ingress.class" = "nginx"
      "ingress.kubernetes.io/rewrite-target" = "/"
    }
  }

  spec {
    rule {
      http {
        path {
          backend {
            service_name = kubernetes_service.test.metadata.0.name
            service_port = 5678
          }

          path = "/test"
        }
      }
    }
  }
}

data "kubernetes_service" "nginx-ingress-controller" {  
  depends_on = [
    helm_release.nginx_ingress
  ]
  metadata {
    name      = "nginx-ingress-controller"
    namespace = kubernetes_namespace.test.metadata.0.name
  }
}
