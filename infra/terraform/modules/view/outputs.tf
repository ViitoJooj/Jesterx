output "service_name" {
  value = kubernetes_service.view.metadata[0].name
}
