variable "kubeconfig_path" {
  description = "Path to kubeconfig used to reach the cluster"
  type        = string
  default     = "~/.kube/config"
}

variable "kube_context" {
  description = "kubeconfig context to use"
  type        = string
  default     = null
}

variable "namespace" {
  description = "Namespace all clusters (view/daemon/storage/caddy) are deployed into"
  type        = string
  default     = "default"
}

variable "env_file_path" {
  description = "Path to the root .env file shared by every cluster"
  type        = string
  default     = "../../.env"
}

variable "view_image" {
  type    = string
  default = "jesterxview:latest"
}

variable "daemon_image" {
  type    = string
  default = "jesterxdaemon:latest"
}

variable "postgres_image" {
  type    = string
  default = "postgres:16-alpine"
}

variable "postgres_db" {
  type    = string
  default = "database"
}

variable "view_replicas" {
  type    = number
  default = 2
}

variable "daemon_replicas" {
  type    = number
  default = 3
}

variable "storage_replicas" {
  type    = number
  default = 3
}

variable "enable_redis" {
  description = "Turn on the redis container in the storage cluster"
  type        = bool
  default     = false
}

variable "redis_image" {
  type    = string
  default = "redis:7-alpine"
}

variable "enable_mongo" {
  description = "Turn on the mongo container in the storage cluster"
  type        = bool
  default     = false
}

variable "mongo_image" {
  type    = string
  default = "mongo:7"
}

variable "caddy_image" {
  type    = string
  default = "caddy:2-alpine"
}

variable "caddyfile_path" {
  description = "Path to the shared Caddyfile"
  type        = string
  default     = "../caddy/Caddyfile"
}
