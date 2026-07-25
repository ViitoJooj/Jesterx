variable "namespace" {
  type = string
}

variable "image" {
  type = string
}

variable "replicas" {
  type = number
}

variable "env_vars" {
  type = map(string)
}
