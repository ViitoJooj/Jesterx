resource "kubernetes_secret" "postgres" {
  metadata {
    name      = "postgres-secret"
    namespace = var.namespace
  }

  data = var.env_vars
}

# ponytail: replicas>1 here means N independent postgres pods, not a
# replicated cluster (no primary/standby). Matches the literal "db com N
# replicas" ask. The headless service lets the app pin writes to postgres-0
# instead of round-robining. Wire up real streaming replication (or an
# operator) when cross-replica consistency actually matters.
resource "kubernetes_stateful_set" "postgres" {
  metadata {
    name      = "postgres"
    namespace = var.namespace
    labels    = { app = "postgres" }
  }

  spec {
    service_name = "postgres"
    replicas     = var.replicas

    selector {
      match_labels = { app = "postgres" }
    }

    template {
      metadata {
        labels = { app = "postgres" }
      }

      spec {
        container {
          name  = "postgres"
          image = var.postgres_image

          port {
            container_port = 5432
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.postgres.metadata[0].name
            }
          }

          env {
            name  = "POSTGRES_DB"
            value = var.postgres_db
          }

          volume_mount {
            name       = "postgres-data"
            mount_path = "/var/lib/postgresql/data"
          }
        }
      }
    }

    volume_claim_template {
      metadata {
        name = "postgres-data"
      }

      spec {
        access_modes = ["ReadWriteOnce"]

        resources {
          requests = {
            storage = "5Gi"
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "postgres" {
  metadata {
    name      = "postgres"
    namespace = var.namespace
  }

  spec {
    cluster_ip = "None"
    selector   = { app = "postgres" }

    port {
      port        = 5432
      target_port = 5432
    }
  }
}

# em breve: redis entra no cluster de storage junto com o postgres.
resource "kubernetes_deployment" "redis" {
  count = var.enable_redis ? 1 : 0

  metadata {
    name      = "redis"
    namespace = var.namespace
    labels    = { app = "redis" }
  }

  spec {
    replicas = 1

    selector {
      match_labels = { app = "redis" }
    }

    template {
      metadata {
        labels = { app = "redis" }
      }

      spec {
        container {
          name  = "redis"
          image = var.redis_image

          port {
            container_port = 6379
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "redis" {
  count = var.enable_redis ? 1 : 0

  metadata {
    name      = "redis"
    namespace = var.namespace
  }

  spec {
    selector = { app = "redis" }

    port {
      port        = 6379
      target_port = 6379
    }
  }
}

# em breve: mongo entra no cluster de storage junto com o postgres.
resource "kubernetes_deployment" "mongo" {
  count = var.enable_mongo ? 1 : 0

  metadata {
    name      = "mongo"
    namespace = var.namespace
    labels    = { app = "mongo" }
  }

  spec {
    replicas = 1

    selector {
      match_labels = { app = "mongo" }
    }

    template {
      metadata {
        labels = { app = "mongo" }
      }

      spec {
        container {
          name  = "mongo"
          image = var.mongo_image

          port {
            container_port = 27017
          }

          env_from {
            secret_ref {
              name = kubernetes_secret.postgres.metadata[0].name
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "mongo" {
  count = var.enable_mongo ? 1 : 0

  metadata {
    name      = "mongo"
    namespace = var.namespace
  }

  spec {
    selector = { app = "mongo" }

    port {
      port        = 27017
      target_port = 27017
    }
  }
}
