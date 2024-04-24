job "datadog-agent-logs-collector" {
  datacenters = ["dc1"]

  group "datadog-agent-group" {
    count = 1

    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }

      port "tcp" {
        to = 10500
      }
    }

    task "datadog-agent" {
      driver = "exec"

      service {
        name = "datadog-agent-logs-collector"
        port = "tcp"
        tags = ["tcp"]
        provider = "nomad"
      }

      template {
        data = <<-EOH
        {{- with nomadVar "nomad/jobs" }}
          REQUESTS_CA_BUNDLE="{{ .MEV_COMMIT_CA_TLS_CERTIFICATE_FILE }}"
        {{ end }}
        EOH
        destination = "local/variables.env"
        env = true
      }

      template {
        data = <<-EOH
          site: datadoghq.com
          logs_enabled: true
        {{- with nomadVar "nomad/jobs" }}
          api_key: {{ .DATA_DOG_API_KEY }}
          tags:
            - env: {{ .DATA_DOG_ENV_TAG }}
              version: {{ .DATA_DOG_VERSION_TAG }}
        {{- end }}
        EOH
        destination = "etc/datadog-agent/datadog.yaml"
      }

      template {
        data = <<-EOH
        {{- with nomadVar "nomad/jobs" }}
          logs:
            - type: tcp
              port: 10500
              source: "{{ .DATA_DOG_ENV_TAG }}-{{ .DATA_DOG_VERSION_TAG }}"
              encoding: "utf-8"
        {{- end }}
        EOH
        destination = "etc/datadog-agent/conf.d/logs-collector.d/conf.yaml"
      }

      config {
        command = "datadog-agent"
        args = ["run"]
      }
    }
  }
}
