job "mev-commit-geth-member-node" {
  datacenters = ["dc1"]

  meta {
    MEV_COMMIT_GETH_VERSION = "v0.3.0-rc1"
  }

  group "member-node-group" {
    count = 1

    scaling {
      min = 1
      max = 1
    }

    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }

      port "metrics" {
        to = 6060
      }
      port "http" {
        static = 8555
        to = 8545
      }
      port "ws" {
        static = 8556
        to = 8546
      }
      port "p2p" {
        to = 30311
      }
    }

    service {
      name = "mev-commit-geth-member-node"
      port = "metrics"
      tags = ["metrics"]
      provider = "nomad"
    }

    task "member-node" {
      driver = "exec"

      service {
        name = "mev-commit-geth-member-node"
        port = "http"
        tags = ["http"]
        provider = "nomad"
      }

      service {
        name = "mev-commit-geth-member-node"
        port = "ws"
        tags = ["ws"]
        provider = "nomad"
      }

      service {
        name = "mev-commit-geth-member-node"
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
          GETH_LOG_TAGS="service:mev-commit-geth-member-node-{{env "NOMAD_ALLOC_INDEX"}}"
          GETH_NODE_TYPE="member"
          GETH_GETH_SYNC_MODE="snap"
          GENESIS_L1_PATH="local/genesis.json"
          GETH_DATA_DIR="local/data-{{env "NOMAD_ALLOC_INDEX"}}"
          GETH_BIN_PATH="local/geth"
        {{- range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "p2p" .Tags }}
          BOOTNODE_ENDPOINT="enode://34a2a388ad31ca37f127bb9ffe93758ee711c5c2277dff6aff2e359bcf2c9509ea55034196788dbd59ed70861f523c1c03d54f1eabb2b4a5c1c129d966fe1e65@{{ .Address }}:{{ .Port }}"
          {{- end }}
        {{- end }}
          NODE_IP="0.0.0.0"
        {{- with nomadVar "nomad/jobs" }}
          {{- $ips := split "," (printf "%s" .MEV_COMMIT_GETH_MEMBER_NODE_IP_ADDRESS_POOL) }}
          {{- $idx := env "NOMAD_ALLOC_INDEX" | parseInt }}
          PUBLIC_NODE_IP={{ index $ips $idx }}
        {{ end }}
          NET_RESTRICT="0.0.0.0/0"
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

          if [ -z "$PUBLIC_NODE_IP" ]; then
            echo "Error: PUBLIC_NODE_IP variable is not set"
            exit 1
          fi

          mkdir -p ${GETH_DATA_DIR} > /dev/null 2>&1
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
