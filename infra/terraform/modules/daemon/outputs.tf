output "service_name" {
  value = kubernetes_service.daemon.metadata[0].name
}
