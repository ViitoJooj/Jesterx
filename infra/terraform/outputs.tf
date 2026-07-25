output "view_service" {
  value = module.view.service_name
}

output "daemon_service" {
  value = module.daemon.service_name
}

output "postgres_service" {
  value = module.storage.postgres_service_name
}

output "caddy_service" {
  value = module.caddy.service_name
}
