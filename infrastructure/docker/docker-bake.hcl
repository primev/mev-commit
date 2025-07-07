variable "TAG" { default = "dev" }

target "mev-commit-builder" {
  context    = "../../"
  dockerfile = "infrastructure/docker/Dockerfile.builder"
}

target "mev-commit-oracle" {
  context    = "./"
  dockerfile = "Dockerfile.oracle"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/mev-commit-oracle:${TAG}"]
}

target "mev-commit" {
  context    = "./"
  dockerfile = "Dockerfile.p2p"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/mev-commit:${TAG}"]
}

target "mev-commit-bridge" {
  context    = "./"
  dockerfile = "Dockerfile.bridge"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/mev-commit-bridge:${TAG}"]
}

target "mev-commit-dashboard" {
  context    = "./"
  dockerfile = "Dockerfile.dashboard"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/mev-commit-dashboard:${TAG}"]
}

target "preconf-rpc" {
  context    = "./"
  dockerfile = "Dockerfile.rpc"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/preconf-rpc:${TAG}"]
}

target "bidder-emulator" {
  context    = "./"
  dockerfile = "Dockerfile.bidderemulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/bidder-emulator:${TAG}"]
}

target "provider-emulator" {
  context    = "./"
  dockerfile = "Dockerfile.provideremulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/provider-emulator:${TAG}"]
}

target "relay-emulator" {
  context    = "./"
  dockerfile = "Dockerfile.relayemulator"
  contexts = {
    builder_ctx = "target:mev-commit-builder"
  }
  tags = ["ghcr.io/primev/relay-emulator:${TAG}"]
}

group "default" {
  targets = ["mev-commit-builder", "mev-commit-oracle", "mev-commit", "mev-commit-bridge", "mev-commit-dashboard", "preconf-rpc", "bidder-emulator", "provider-emulator", "relay-emulator"]
}

