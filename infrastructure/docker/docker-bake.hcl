variable "TAG" { default = "dev" }
variable "PLATFORM" { default = ["linux/amd64"] }
variable "REGISTRY" { default = "ghcr.io/primev" }
variable "REPO_NAME" { default = "" }

# Git variables - these will be passed from Makefile
variable "GIT_BRANCH" { 
  default = null
}

variable "GIT_COMMIT" { 
  default = null
}

function "get_labels" {
  params = [component]
  result = {
    "branch" = GIT_BRANCH != null ? GIT_BRANCH : "unknown"
    "commit" = GIT_COMMIT != null ? GIT_COMMIT : "unknown"
    "component" = component
    "build.timestamp" = timestamp()
    "build.tag" = TAG
  }
}

target "mev-commit-builder" {
  inherits   = ["_common"]
  context    = "../../"
  dockerfile = "infrastructure/docker/Dockerfile.builder"
  labels = get_labels("builder")
}

target "mev-commit-oracle" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.oracle"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit-oracle" : "${REGISTRY}/mev-commit-oracle:${TAG}"]
  labels = get_labels("oracle")
}

target "mev-commit" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.p2p"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit" : "${REGISTRY}/mev-commit:${TAG}"]
  labels = get_labels("p2p")
}

target "mev-commit-bridge" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.bridge"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit-bridge" : "${REGISTRY}/mev-commit-bridge:${TAG}"]
  labels = get_labels("bridge")
}

target "mev-commit-dashboard" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.dashboard"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-mev-commit-dashboard" : "${REGISTRY}/mev-commit-dashboard:${TAG}"]
  labels = get_labels("dashboard")
}

target "preconf-rpc" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.rpc"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-preconf-rpc" : "${REGISTRY}/preconf-rpc:${TAG}"]
  labels = get_labels("preconf-rpc")
}

target "bidder-emulator" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.bidderemulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-bidder-emulator" : "${REGISTRY}/bidder-emulator:${TAG}"]
  labels = get_labels("bidder-emulator")
}

target "provider-emulator" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.provideremulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-provider-emulator" : "${REGISTRY}/provider-emulator:${TAG}"]
  labels = get_labels("provider-emulator")
}

target "relay-emulator" {
  inherits = ["_common"]
  context    = "./"
  dockerfile = "Dockerfile.relayemulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = [REPO_NAME != "" ? "${REGISTRY}/${REPO_NAME}:${TAG}-relay-emulator" : "${REGISTRY}/relay-emulator:${TAG}"]
  labels = get_labels("relay-emulator")
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
