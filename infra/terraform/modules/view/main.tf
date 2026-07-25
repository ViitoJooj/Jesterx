resource "kubernetes_secret" "view" {
  metadata {
    name      = "view-secret"
    namespace = var.namespace
  }

  data = var.env_vars
}

resource "kubernetes_deployment" "view" {
  metadata {
    name      = "view"
    namespace = var.namespace
    labels    = { app = "view" }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = { app = "view" }
    }

    template {
      metadata {
        labels = { app = "view" }
      }

      spec {
        container {
          name  = "jesterx-view"
          image = var.image

          port {
            container_port = 80
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.view.metadata[0].name
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "view" {
  metadata {
    name      = "view"
    namespace = var.namespace
  }

  spec {
    selector = { app = "view" }

    port {
      port        = 80
      target_port = 80
    }
  }
}
