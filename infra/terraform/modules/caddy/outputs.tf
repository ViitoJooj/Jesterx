output "service_name" {
  value = kubernetes_service.caddy.metadata[0].name
}
