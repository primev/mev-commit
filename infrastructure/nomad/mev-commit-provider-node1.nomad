job "mev-commit-provider-node1" {
  datacenters = ["dc1"]

  meta {
    MEV_COMMIT_VERSION = "v0.3.0-rc2"
    KEYSTORE_GENERATOR_VERSION_ID = "RYncrfSOWyiyojCLlse4f4kSiotrdhgM"
  }

  group "provider-node-group" {
    count = 1

    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }

      port "metrics" {
        to = 13523
      }
      port "http" {
        static = 13533
        to = 13523
      }
      port "p2p" {
        to = 13522
      }
      port "rpc" {
        static = 13534
        to = 13524
      }
    }

    service {
      name = "mev-commit-provider-node1"
      port = "metrics"
      tags = ["metrics"]
      provider = "nomad"
    }

    service {
      name = "mev-commit-provider-node1"
      port = "p2p"
      tags = ["p2p"]
      provider = "nomad"
    }

    service {
      name = "mev-commit-provider-node1"
      port = "http"
      tags = ["http"]
      provider = "nomad"
    }

    service {
      name = "mev-commit-provider-node1"
      port = "rpc"
      tags = ["rpc"]
      provider = "nomad"
    }

    task "provider-node" {
      driver = "exec"

      artifact {
        source = "https://github.com/primevprotocol/mev-commit/releases/download/${NOMAD_META_MEV_COMMIT_VERSION}/mev-commit_Linux_x86_64.tar.gz"
      }

      artifact {
        source = "https://primev-infrastructure-artifacts.s3.us-west-2.amazonaws.com/keystore-generator?versionId=${NOMAD_META_KEYSTORE_GENERATOR_VERSION_ID}"
      }

      template {
        data = <<-EOH
          KEYSTOREGEN_LOG_FMT="json"
          KEYSTOREGEN_LOG_TAGS="service:mev-commit-provider-node1"
          MEV_COMMIT_LOG_FMT="json"
          MEV_COMMIT_LOG_TAGS="service:mev-commit-provider-node1"
          MEV_COMMIT_KEYSTORE_PATH = "/local/data-{{env "NOMAD_ALLOC_INDEX"}}/keystore"
          MEV_COMMIT_KEYSTORE_PASSWORD = "{{ with secret "secret/data/mev-commit" }}{{ .Data.data.provider1_keystore_password }}{{ end }}"
      {{- range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "http" .Tags }}
          MEV_COMMIT_SETTLEMENT_RPC_ENDPOINT="http://{{ .Address }}:{{ .Port }}"
          {{- end }}
      {{- end }}
          MEV_COMMIT_PEER_TYPE = "provider"
      {{- with nomadVar "nomad/jobs" }}
          MEV_COMMIT_NAT_ADDR = "{{ .MEV_COMMIT_PROVIDER_NODE1_NAT_ADDR }}"
      {{ end }}
      {{- range nomadService "mev-commit-provider-node1" }}
          {{- if contains "p2p" .Tags }}
          MEV_COMMIT_NAT_PORT="{{ .Port }}"
          {{- end }}
      {{- end }}
          MEV_COMMIT_HTTP_ADDR = "0.0.0.0"
          MEV_COMMIT_RPC_ADDR = "0.0.0.0"
          MEV_COMMIT_P2P_ADDR = "0.0.0.0"
        {{- with nomadVar "nomad/jobs" }}
          MEV_COMMIT_SERVER_TLS_CERTIFICATE="{{ .MEV_COMMIT_SERVER_TLS_CERTIFICATE_FILE }}"
          MEV_COMMIT_SERVER_TLS_PRIVATE_KEY="{{ .MEV_COMMIT_SERVER_TLS_PRIVATE_KEY_FILE }}"
        {{ end }}
        {{- range nomadService "datadog-agent-logs-collector" }}
          {{- if contains "mev-commit-provider-node1-sink" .Tags }}
          MEV_COMMIT_LOG_SINK_TCP="{{ .Address }}:{{ .Port }}"
          {{ end }}
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

        {{- range nomadService "mev-commit-bootnode1" }}
          {{- if contains "http" .Tags }}

          # Fetch topology data
          TOPOLOGY=$(curl https://{{ .Address }}:{{ .Port }}/topology)

          # Extract peerId (assuming the response format is JSON and jq is available)
          PEER_ID=$(echo ${TOPOLOGY} | jq -r '.self.Underlay')

          {{- else if contains "p2p" .Tags }}

          # Set the MEV_COMMIT_BOOTNODES environment variable
          export MEV_COMMIT_BOOTNODES="/ip4/{{ .Address }}/tcp/{{ .Port }}/p2p/${PEER_ID}"

          {{- end }}
        {{- end }}

          if [ ! -d "${MEV_COMMIT_KEYSTORE_PATH}" ]; then
            mkdir -p "${MEV_COMMIT_KEYSTORE_PATH}" > /dev/null 2>&1
            chmod +x local/keystore-generator
            local/keystore-generator generate \
              --keystore-dir "${MEV_COMMIT_KEYSTORE_PATH}" \
              --passphrase "${MEV_COMMIT_KEYSTORE_PASSWORD}"
          fi

          chmod +x local/mev-commit
          local/mev-commit
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
