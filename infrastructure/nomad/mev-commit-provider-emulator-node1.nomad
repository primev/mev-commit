job "mev-commit-provider-emulator-node1" {
  datacenters = ["dc1"]

  meta {
    MEV_COMMIT_PROVIDER_EMULATOR_VERSION_ID = "Kptl3SR.e75jNaO4YczMA.vIbabj7A.u"
  }

  group "provider-emulator-node-group" {
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

    task "provider-emulator-node" {
      driver = "exec"

      service {
        name = "mev-commit-provider-emulator-node1"
        port = "metrics"
        tags = ["metrics"]
        provider = "nomad"
      }

      artifact {
        source = "https://primev-infrastructure-artifacts.s3.us-west-2.amazonaws.com/provider-emulator?versionId=${NOMAD_META_MEV_COMMIT_PROVIDER_EMULATOR_VERSION_ID}"
      }

      template {
        data = <<-EOH
          #!/usr/bin/env bash

        {{- range nomadService "datadog-agent-logs-collector" }}
          {{ if contains "tcp" .Tags }}
          exec > >(nc {{ .Address }} {{ .Port }}) 2>&1
          {{ end }}
        {{- end }}

        {{- range nomadService "mev-commit-provider-node1" }}
          {{- if contains "rpc" .Tags }}
          PROVIDER_IP_PORT="{{ .Address }}:{{ .Port }}"
          {{- end }}
        {{- end }}

          LOG_TAGS="service:mev-commit-provider-emulator-node1"

          chmod +x local/provider-emulator
          local/provider-emulator \
            -server-addr "${PROVIDER_IP_PORT}" \
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
