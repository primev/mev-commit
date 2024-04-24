job "mev-commit-bridge" {
  datacenters = ["dc1"]

  meta {
    CONTRACTS_VERSION = "v0.2.0"
    MEV_COMMIT_BRIDGE_VERSION = "v0.0.8"
  }

  group "bridge-group" {
    restart {
      attempts = 0
      interval = "30m"
      delay = "15s"
      mode = "fail"
    }
    count = 1

    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }
      
      port "metrics" {
        to = 8080
      }
    }

    task "deploy-and-relay" {
      driver = "exec"

      service {
        name = "mev-commit-bridge"
        port = "metrics"
        tags = ["metrics"]
        provider = "nomad"
      }

      artifact {
        source = "https://github.com/foundry-rs/foundry/releases/download/nightly-293fad73670b7b59ca901c7f2105bf7a29165a90/foundry_nightly_linux_amd64.tar.gz"
      }

      artifact {
        source = "git::https://github.com/primevprotocol/contracts"
        destination = "local/contracts"
        options {
          ref = "${NOMAD_META_CONTRACTS_VERSION}"
          depth = 1
        }
      }

      artifact {
        source = "https://raw.githubusercontent.com/primevprotocol/mev-commit-bridge/${NOMAD_META_MEV_COMMIT_BRIDGE_VERSION}/standard/bridge-v1/deploy_contracts.sh"
      }

      artifact {
        source = "https://github.com/primevprotocol/mev-commit-bridge/releases/download/${NOMAD_META_MEV_COMMIT_BRIDGE_VERSION}/mev-commit-bridge-relayer-linux-amd64.tar.gz"
      }

      template {
        data = "eaa870cc825f07a18b1cfe7bb62a03c6b2601e1129730c9a8724a2df6eeea4f4"
        destination = "local/relayer_key"
        perms = "0600"
      }

      template {
        data = <<-EOH
          STANDARD_BRIDGE_RELAYER_LOG_FMT="json"
          STANDARD_BRIDGE_RELAYER_LOG_LEVEL="debug"
          STANDARD_BRIDGE_RELAYER_LOG_TAGS="service:mev-commit-bridge"
        {{- range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "http" .Tags }}
          STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL="http://{{ .Address }}:{{ .Port }}"
          {{- end }}
        {{- end }}
          STANDARD_BRIDGE_RELAYER_L1_RPC_URL="https://ethereum-holesky.publicnode.com"
          STANDARD_BRIDGE_RELAYER_PRIV_KEY_FILE="local/relayer_key"
          L1_CHAIN_ID=17000
          RELAYER_PRIVKEY=0xeaa870cc825f07a18b1cfe7bb62a03c6b2601e1129730c9a8724a2df6eeea4f4
          FORGE_BIN_PATH="local/forge"
          CAST_BIN_PATH="local/cast"
          CONTRACTS_PATH="local/contracts"
          ARTIFACT_OUT_PATH="local"
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

        if [ -f "local/L1GatewayArtifact.json" ] && [ -f "local/SettlementGatewayArtifact.json" ]; then
          echo "Artifacts exist. Skipping contract deployment..."
        else
          echo "Deploying contracts..."
          chmod +x local/deploy_contracts.sh
          ./local/deploy_contracts.sh
        fi

        export STANDARD_BRIDGE_RELAYER_L1_CONTRACT_ADDR="$(jq -r '.l1_gateway_addr' local/L1GatewayArtifact.json)"
        export STANDARD_BRIDGE_RELAYER_SETTLEMENT_CONTRACT_ADDR="$(jq -r '.settlement_gateway_addr' local/SettlementGatewayArtifact.json)"

        if [ -z "$STANDARD_BRIDGE_RELAYER_L1_CONTRACT_ADDR" ] || ! echo "$STANDARD_BRIDGE_RELAYER_L1_CONTRACT_ADDR" | grep -q "^0x"; then
          echo "Error: L1 Contract Address is not populated. Exiting.."
          exit 1
        fi
        if [ -z "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_CONTRACT_ADDR" ] || ! echo "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_CONTRACT_ADDR" | grep -q "^0x"; then
            echo "Error: Settlement Contract Address is not populated. Exiting.."
            exit 1
        fi

        echo "L1 Contract Address: $STANDARD_BRIDGE_RELAYER_L1_CONTRACT_ADDR"
        echo "Settlement Contract Address: $STANDARD_BRIDGE_RELAYER_SETTLEMENT_CONTRACT_ADDR"

        chmod +x local/relayer-linux-amd64
        ./local/relayer-linux-amd64 start
        EOH
        destination = "local/run.sh"
        perms = "0755"
      }

      config {
        command = "bash"
        args    = ["-c", "local/run.sh"]
      }
    }
  }
}
