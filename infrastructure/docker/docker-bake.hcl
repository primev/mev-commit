variable "TAG" { default = "dev" }
variable "PLATFORM" { default = ["linux/amd64"] }
variable "REGISTRY" { default = "ghcr.io/primev" }
variable "REPO_NAME" { default = "" }

target "mev-commit-builder" {
  inherits   = ["_common"]
  context    = "../../"
  dockerfile = "infrastructure/docker/Dockerfile.builder"
}

target "mev-commit-oracle" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.oracle"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit-oracle" : "${REGISTRY}/mev-commit-oracle:${TAG}"]
}

target "mev-commit" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.p2p"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit" : "${REGISTRY}/mev-commit:${TAG}"]
}

target "mev-commit-bridge" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.bridge"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit-bridge" : "${REGISTRY}/mev-commit-bridge:${TAG}"]
}

target "mev-commit-dashboard" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.dashboard"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit-dashboard" : "${REGISTRY}/mev-commit-dashboard:${TAG}"]
}

target "preconf-rpc" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.rpc"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-preconf-rpc" : "${REGISTRY}/preconf-rpc:${TAG}"]
}

target "bidder-emulator" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.bidderemulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-bidder-emulator" : "${REGISTRY}/bidder-emulator:${TAG}"]
}

target "provider-emulator" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.provideremulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-provider-emulator" : "${REGISTRY}/provider-emulator:${TAG}"]
}

target "relay-emulator" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.relayemulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-relay-emulator" : "${REGISTRY}/relay-emulator:${TAG}"]
}

group "all" {
  targets = ["mev-commit-builder", "mev-commit-oracle", "mev-commit", "mev-commit-bridge", "mev-commit-dashboard", "preconf-rpc", "bidder-emulator", "provider-emulator", "relay-emulator"]
}

group "default" {
  targets = ["all"]
}

target "_common" {
  platforms = PLATFORM
  output = ["type=docker"]
}
