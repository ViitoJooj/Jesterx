module "view" {
  source    = "./modules/view"
  namespace = var.namespace
  image     = var.view_image
  replicas  = var.view_replicas
  env_vars  = local.env_vars
}

module "daemon" {
  source    = "./modules/daemon"
  namespace = var.namespace
  image     = var.daemon_image
  replicas  = var.daemon_replicas
  env_vars  = local.env_vars
}

module "storage" {
  source         = "./modules/storage"
  namespace      = var.namespace
  postgres_image = var.postgres_image
  postgres_db    = var.postgres_db
  replicas       = var.storage_replicas
  env_vars       = local.env_vars

  enable_redis = var.enable_redis
  redis_image  = var.redis_image

  enable_mongo = var.enable_mongo
  mongo_image  = var.mongo_image
}

module "caddy" {
  source         = "./modules/caddy"
  namespace      = var.namespace
  image          = var.caddy_image
  caddyfile_path = var.caddyfile_path
}
