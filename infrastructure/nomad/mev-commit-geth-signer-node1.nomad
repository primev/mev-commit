job "mev-commit-geth-signer-node1" {
  datacenters = ["dc1"]

  meta {
    MEV_COMMIT_GETH_VERSION = "v0.3.0-rc1"
  }

  group "signer-node-group" {
    count = 1

    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }
 
      port "metrics" {
        to = 6060
      }
      port "p2p" {
        to = 30311
      }
    }

    task "signer-node" {
      driver = "exec"

      service {
        name = "mev-commit-geth-signer-node1"
        port = "metrics"
        tags = ["metrics"]
        provider = "nomad"
      }

      service {
        name = "mev-commit-geth-signer-node1"
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

      artifact {
        source = "https://raw.githubusercontent.com/primevprotocol/mev-commit-geth/${NOMAD_META_MEV_COMMIT_GETH_VERSION}/geth-poa/signer-node1/keystore/UTC--2024-02-19T14-30-04.415601156Z--d9cd8e5de6d55f796d980b818d350c0746c25b97"
        destination = "local/data-${NOMAD_ALLOC_INDEX}/keystore"
      }

      template {
        data = <<-EOH
          GETH_LOG_FORMAT="json"
          GETH_LOG_TAGS="service:mev-commit-geth-signer-node1"
          GETH_NODE_TYPE = "signer"
          GENESIS_L1_PATH = "local/genesis.json"
          GETH_DATA_DIR = "local/data-{{env "NOMAD_ALLOC_INDEX"}}"
          MEV_COMMIT_GETH_PASSWORD = "{{ with secret "secret/data/mev-commit" }}{{ .Data.data.geth_signer1_keystore_password }}{{ end }}"
          GETH_BIN_PATH = "local/geth"
          BLOCK_SIGNER_ADDRESS = "0xd9cd8E5DE6d55f796D980B818D350C0746C25b97"
          NODE_IP = "0.0.0.0"
          NET_RESTRICT= "0.0.0.0/0"
        {{- range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "p2p" .Tags }}
          BOOTNODE_ENDPOINT="enode://34a2a388ad31ca37f127bb9ffe93758ee711c5c2277dff6aff2e359bcf2c9509ea55034196788dbd59ed70861f523c1c03d54f1eabb2b4a5c1c129d966fe1e65@{{ .Address }}:{{ .Port }}"
          {{- end }}
        {{- end }}
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

