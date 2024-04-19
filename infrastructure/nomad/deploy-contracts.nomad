job "contract-deployment" {
  datacenters = ["dc1"]
  type = "batch"

  meta {
    CONTRACTS_VERSION = "v0.2.0"
    MEV_COMMIT_GETH_VERSION = "v0.3.0-rc1"
  }

  group "contract-deployment-group" {
    restart {
      attempts = 0
      interval = "30m"
      delay = "15s"
      mode = "fail"
    }
    count = 1

    network {
      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }
    }

    task "deploy-create2" {
      driver = "exec"

      lifecycle {
        hook = "prestart"
      }

      artifact {
        source = "https://raw.githubusercontent.com/primevprotocol/mev-commit-geth/${NOMAD_META_MEV_COMMIT_GETH_VERSION}/geth-poa/util/deploy_create2.sh"
      }

      template {
        data = <<-EOH
          #!/usr/bin/env bash

        {{- range nomadService "datadog-agent-logs-collector" }}
          {{ if contains "tcp" .Tags }}
          exec > >(nc {{ .Address }} {{ .Port }}) 2>&1
          {{ end }}
        {{- end }}

        {{- range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "http" .Tags }}
          chmod +x local/deploy_create2.sh
          local/deploy_create2.sh "http://{{ .Address }}:{{ .Port }}"
          {{- end }}
        {{- end }}
        EOH
        destination = "local/run.sh"
        perms = "0755"
      }

      config {
        command = "bash"
        args = ["-c", "local/run.sh"]
      }
    }

    task "deploy-contracts" {
      driver = "exec"

      artifact {
        source = "https://github.com/foundry-rs/foundry/releases/download/nightly-293fad73670b7b59ca901c7f2105bf7a29165a90/foundry_nightly_linux_amd64.tar.gz"
      }

      artifact {
        source      = "git::https://github.com/primevprotocol/contracts.git"
        destination = "local/contracts"
        options {
          ref = "${NOMAD_META_CONTRACTS_VERSION}"
          depth = 1
        }
      }

      template {
        data = <<-EOH
        {{- range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "http" .Tags }}
          RPC_URL = "http://{{ .Address }}:{{ .Port }}"
          {{- end }}
        {{- end }}
          PRIVATE_KEY = "{{ with secret "secret/data/mev-commit" }}{{ .Data.data.deploy_contracts_private_key }}{{ end }}"
          CHAIN_ID = "17864"
          FORGE_BIN_PATH = "local/forge"
          DEPLOY_TYPE = "core"
          SCRIPT_PATH_PREFIX = "local/contracts/scripts/"
          CONTRACT_REPO_ROOT_PATH = "local/contracts"
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

          local/forge --version && local/forge build --use 0.8.23 --root local/contracts
          chmod +x local/contracts/entrypoint.sh && local/contracts/entrypoint.sh
          rm -rf local/contracts
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
