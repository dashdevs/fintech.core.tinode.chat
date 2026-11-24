# Builds the Tinode image from this fork's sources on top of the unmodified upstream
# docker/tinode/Dockerfile.
#
# $ TARGET_DB=postgres docker buildx bake            # build tinode for both platforms
# $ TARGET_DB=mysql docker buildx bake --load \
#     --set tinode.platform=linux/arm64 --set upstream.platform=linux/arm64

# Database backend: postgres, mysql, mongodb, rethinkdb or alldbs.
variable "TARGET_DB" { default = "postgres" }

# Upstream release whose archive provides the runtime image and the web client.
variable "UPSTREAM_VERSION" { default = "0.25.3" }

# Version of this fork, written to the VERSION env and the image label.
variable "VERSION" { default = UPSTREAM_VERSION }

# Value of main.buildstamp in the server binary.
variable "BUILDSTAMP" { default = "unknown" }

group "default" {
  targets = ["tinode"]
}

target "_platforms" {
  platforms = ["linux/amd64", "linux/arm64"]
}

# Filled in by docker/metadata-action in CI (tags, labels); empty locally.
target "docker-metadata-action" {}

target "upstream" {
  inherits   = ["_platforms"]
  context    = "docker/tinode"
  dockerfile = "Dockerfile"
  args = {
    TARGET_DB = TARGET_DB
    VERSION   = UPSTREAM_VERSION
  }
}

target "tinode" {
  inherits   = ["_platforms", "docker-metadata-action"]
  context    = "."
  dockerfile = "docker/tinode/Dockerfile.source"
  contexts = {
    base = "target:upstream"
  }
  args = {
    TARGET_DB   = TARGET_DB
    APP_VERSION = VERSION
    BUILDSTAMP  = BUILDSTAMP
  }
}
