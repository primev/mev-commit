job "mev-commit-bidder-emulator-node1" {
  datacenters = ["dc1"]

  meta {
    MEV_COMMIT_BIDDER_EMULATOR_VERSION_ID = "aSC397vNgaBKLH6Hx8ckrdiUKcFf7pDW"
  }

  group "bidder-emulator-node-group" {
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

    task "bidder-emulator-node" {
      driver = "exec"

      service {
        name = "mev-commit-bidder-emulator-node1"
        port = "metrics"
        tags = ["metrics"]
        provider = "nomad"
      }

      artifact {
        source = "https://primev-infrastructure-artifacts.s3.us-west-2.amazonaws.com/bidder-emulator?versionId=${NOMAD_META_MEV_COMMIT_BIDDER_EMULATOR_VERSION_ID}"
      }

      template {
        data = <<-EOH
          #!/usr/bin/env bash

        {{- range nomadService "datadog-agent-logs-collector" }}
          {{ if contains "tcp" .Tags }}
          exec > >(nc {{ .Address }} {{ .Port }}) 2>&1
          {{ end }}
        {{- end }}

        {{ range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "http" .Tags }}
          GETH_RPC_ADDRESS="http://{{ .Address }}:{{ .Port }}"
          {{ end }}
        {{ end }}

        {{- range nomadService "mev-commit-bidder-node1" }}
          {{- if contains "rpc" .Tags }}
          BIDDER_IP_PORT="{{ .Address }}:{{ .Port }}"
          {{- end }}
        {{- end }}

          LOG_TAGS="service:mev-commit-bidder-emulator-node1"

          chmod +x local/bidder-emulator
          local/bidder-emulator \
            -server-addr "${BIDDER_IP_PORT}" \
            -rpc-addr "${GETH_RPC_ADDRESS}" \
            -log-tags "${LOG_TAGS}" \
            -log-fmt "json"
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
