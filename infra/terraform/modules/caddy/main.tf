resource "kubernetes_config_map" "caddy" {
  metadata {
    name      = "caddy-config"
    namespace = var.namespace
  }

  data = {
    Caddyfile = file(var.caddyfile_path)
  }
}

resource "kubernetes_deployment" "caddy" {
  metadata {
    name      = "caddy"
    namespace = var.namespace
    labels    = { app = "caddy" }
  }

  spec {
    replicas = 1

    selector {
      match_labels = { app = "caddy" }
    }

    template {
      metadata {
        labels = { app = "caddy" }
      }

      spec {
        container {
          name  = "caddy"
          image = var.image

          port {
            container_port = 80
          }

          port {
            container_port = 443
          }

          volume_mount {
            name       = "caddy-config"
            mount_path = "/etc/caddy"
          }
        }

        volume {
          name = "caddy-config"

          config_map {
            name = kubernetes_config_map.caddy.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "caddy" {
  metadata {
    name      = "caddy"
    namespace = var.namespace
  }

  spec {
    type     = "LoadBalancer"
    selector = { app = "caddy" }

    port {
      name        = "http"
      port        = 80
      target_port = 80
    }

    port {
      name        = "https"
      port        = 443
      target_port = 443
    }
  }
}
