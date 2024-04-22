job "datadog-agent-metrics-collector" {
  datacenters = ["dc1"]

  group "datadog-agent-group" {
    count = 1

    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }
    }

    task "datadog-agent" {
      driver = "exec"

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
          init_config:

          instances:
          {{- range nomadService "mev-commit-bootnode1" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: https://{{ .Address }}:{{ .Port }}/metrics
              service: "mev-commit-bootnode1"
              metrics:
                - mev_commit*
                - go*
                - libp2p*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-provider-node1" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: https://{{ .Address }}:{{ .Port }}/metrics
              service: "mev-commit-provider-node1"
              metrics:
                - mev_commit*
                - go*
                - libp2p*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-bidder-node1" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: https://{{ .Address }}:{{ .Port }}/metrics
              service: "mev-commit-bidder-node1"
              metrics:
                - mev_commit*
                - go*
                - libp2p*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-geth-bootnode1" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: http://{{ .Address }}:{{ .Port }}/debug/metrics/prometheus
              service: "mev-commit-geth-bootnode1"
              metrics:
                - txpool*
                - trie*
                - system*
                - state*
                - rpc*
                - p2p*
                - eth*
                - chain*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-geth-member-node" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: http://{{ .Address }}:{{ .Port }}/debug/metrics/prometheus
              service: "mev-commit-geth-member-node"
              metrics:
                - txpool*
                - trie*
                - system*
                - state*
                - rpc*
                - p2p*
                - eth*
                - chain*
                - clique*
                - vm*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-geth-signer-node1" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: http://{{ .Address }}:{{ .Port }}/debug/metrics/prometheus
              service: "mev-commit-geth-signer-node1"
              metrics:
                - txpool*
                - trie*
                - system*
                - state*
                - rpc*
                - p2p*
                - eth*
                - chain*
                - clique*
                - vm*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-geth-signer-node2" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: http://{{ .Address }}:{{ .Port }}/debug/metrics/prometheus
              service: "mev-commit-geth-signer-node2"
              metrics:
                - txpool*
                - trie*
                - system*
                - state*
                - rpc*
                - p2p*
                - eth*
                - chain*
                - clique*
                - vm*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-geth-signer-node3" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: http://{{ .Address }}:{{ .Port }}/debug/metrics/prometheus
              service: "mev-commit-geth-signer-node3"
              metrics:
                - txpool*
                - trie*
                - system*
                - state*
                - rpc*
                - p2p*
                - eth*
                - chain*
                - clique*
                - vm*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-provider-emulator-node1" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: http://{{ .Address }}:{{ .Port }}/metrics
              service: "mev-commit-provider-emulator-node1"
              metrics:
                - mev_commit*
            {{- end }}
          {{- end }}

          {{- range nomadService "mev-commit-bidder-emulator-node1" }}
            {{- if contains "metrics" .Tags }}
            - openmetrics_endpoint: http://{{ .Address }}:{{ .Port }}/metrics
              service: "mev-commit-bidder-emulator-node1"
              metrics:
                - mev_commit*
            {{- end }}
          {{- end }}
        EOH
        destination = "etc/datadog-agent/conf.d/openmetrics.d/conf.yaml"
      }

      config {
        command = "datadog-agent"
        args = ["run"]
      }
    }
  }
}
