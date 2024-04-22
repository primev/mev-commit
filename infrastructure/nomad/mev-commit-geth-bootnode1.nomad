job "mev-commit-geth-bootnode1" {
  datacenters = ["dc1"]

  meta {
    MEV_COMMIT_GETH_VERSION = "v0.3.0-rc1"
  }

  group "bootnode-group" {
    count = 1
    
    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }

      port "metrics" {
        to = 6060
      }
      port "http" {
        static = 8545
        to = 8545
      }
      port "p2p" {
        to = 30301
      }
    }

    task "bootnode" {
      driver = "exec"

      service {
        name = "mev-commit-geth-bootnode1"
        port = "metrics"
        tags = ["metrics"]
        provider = "nomad"
      }

      service {
        name = "mev-commit-geth-bootnode1"
        port = "http"
        tags = ["http"]
        provider = "nomad"
      }

      service {
        name = "mev-commit-geth-bootnode1"
        port = "p2p"
        tags = ["p2p"]
        provider = "nomad"
      }

      artifact {
        source = "https://raw.githubusercontent.com/primevprotocol/mev-commit-geth/${NOMAD_META_MEV_COMMIT_GETH_VERSION}/geth-poa/entrypoint.sh"
        destination = "local/bootstrap-geth.sh"
        mode = "file"
      }

      artifact {
        source = "https://raw.githubusercontent.com/primevprotocol/mev-commit-geth/${NOMAD_META_MEV_COMMIT_GETH_VERSION}/geth-poa/genesis.json"
      }

      artifact {
        source = "https://github.com/primevprotocol/mev-commit-geth/releases/download/${NOMAD_META_MEV_COMMIT_GETH_VERSION}/mev-commit-geth_Linux_x86_64.tar.gz"
      }

      template {
        data = <<-EOH
          GETH_LOG_FORMAT="json"
          GETH_LOG_TAGS="service:mev-commit-geth-bootnode1"
          GETH_NODE_TYPE = "bootnode"
          GENESIS_L1_PATH = "local/genesis.json"
          GETH_DATA_DIR = "local/data-{{env "NOMAD_ALLOC_INDEX"}}"
          GETH_BIN_PATH = "local/geth"
          BOOT_KEY = "7b548c1c0fbe80ef1eb0aaec2edf26fd20fb0d758e94948cf6c5f2a486e735f6"
          NODE_IP = "0.0.0.0"
          PUBLIC_NODE_IP = "0.0.0.0"
          NET_RESTRICT= "0.0.0.0/0"
        EOH
        destination = "local/variables.env"
        env = true
      }

      template {
        data = <<-EOH
          #!/usr/bin/env bash

        {{- range nomadService "datadog-agent-logs-collector" }}
          {{ if contains "tcp" .Tags }}
          exec > >(nc {{ .Address }} {{ .Port }}) 2>&1
          {{ end }}
        {{- end }}

          mkdir -p ${GETH_DATA_DIR}
          chmod +x local/geth local/bootstrap-geth.sh
          local/bootstrap-geth.sh
        EOH
        destination = "local/run.sh"
        perms = "0755"
      }

      config {
        command = "bash"
        args = ["-c", "local/run.sh"]
      }
    }
  }
}

