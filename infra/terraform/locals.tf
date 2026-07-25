locals {
  # parses the root .env into a map so every cluster's Secret comes from the
  # same single source instead of hand-copied values.
  env_lines = [
    for line in split("\n", file(var.env_file_path)) :
    trimspace(line)
    if trimspace(line) != "" && !startswith(trimspace(line), "#")
  ]

  env_vars = {
    for line in local.env_lines :
    trimspace(split("=", line)[0]) => trimspace(join("=", slice(split("=", line), 1, length(split("=", line)))))
  }
}
