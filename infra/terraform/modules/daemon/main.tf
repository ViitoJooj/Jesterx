resource "kubernetes_secret" "daemon" {
  metadata {
    name      = "postgres-secret"
    namespace = var.namespace
  }

  data = var.env_vars
}

resource "kubernetes_deployment" "daemon" {
  metadata {
    name      = "verkoupe"
    namespace = var.namespace
    labels    = { app = "verkoupe" }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = { app = "verkoupe" }
    }

    template {
      metadata {
        labels = { app = "verkoupe" }
      }

      spec {
        container {
          name  = "verkoupe-core"
          image = var.image

          port {
            container_port = 8080
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.daemon.metadata[0].name
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "daemon" {
  metadata {
    name      = "daemon"
    namespace = var.namespace
  }

  spec {
    selector = { app = "verkoupe" }

    port {
      port        = 8080
      target_port = 8080
    }
  }
}
