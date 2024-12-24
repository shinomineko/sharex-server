variable "LABELS" {
  default = {
    "org.opencontainers.image.url" = "https://github.com/shinomineko/sharex-server"
    "org.opencontainers.image.source" = "https://github.com/shinomineko/sharex-server"
    "org.opencontainers.image.created" = timestamp()
    "org.opencontainers.image.revision" = "main"
    "org.opencontainers.image.title" = "sharex-server"
  }
}

target "default" {
  context = "."
  platforms = [ "linux/amd64", "linux/arm64" ]
  labels = LABELS
  tags = [
    "docker.io/shinomineko/sharex-server:main",
    "ghcr.io/shinomineko/sharex-server:main"
  ]
}
