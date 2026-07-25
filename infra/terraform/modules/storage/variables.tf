variable "namespace" {
  type = string
}

variable "postgres_image" {
  type = string
}

variable "postgres_db" {
  type = string
}

variable "replicas" {
  type = number
}

variable "env_vars" {
  type = map(string)
}

variable "enable_redis" {
  type = bool
}

variable "redis_image" {
  type = string
}

variable "enable_mongo" {
  type = bool
}

variable "mongo_image" {
  type = string
}
