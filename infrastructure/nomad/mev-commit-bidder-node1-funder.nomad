job "mev-commit-bidder-node1-funder" {
  datacenters = ["dc1"]
  type = "batch"

  group "mev-commit-bidder-node-funder-group" {
    count = 1

    task "bidder-node-funder" {
      driver = "exec"

      artifact {
        source = "https://github.com/foundry-rs/foundry/releases/download/nightly-293fad73670b7b59ca901c7f2105bf7a29165a90/foundry_nightly_linux_amd64.tar.gz"
      }

      template {
        data = <<-EOH
          #!/usr/bin/env bash

        {{- range nomadService "datadog-agent-logs-collector" }}
          {{ if contains "tcp" .Tags }}
          exec > >(nc {{ .Address }} {{ .Port }}) 2>&1
          {{ end }}
        {{- end }}

        {{- range nomadService "mev-commit-bidder-node1" }}
          {{- if contains "http" .Tags }}
          START_TIME=$(date +%s)
          TIMEOUT=60

          while : ; do
              STATUS_CODE=$(curl -o /dev/null -s -w "%%{http_code}\n" https://{{ .Address }}:{{ .Port }}/health)
              if [ "${STATUS_CODE}" -eq 200 ]; then
                  break
              else
                  CURRENT_TIME=$(date +%s)
                  ELAPSED_TIME=$((CURRENT_TIME - START_TIME))

                  if [ "${ELAPSED_TIME}" -ge "${TIMEOUT}" ]; then
                      echo "Timeout reached; the bidder node is not up."
                      exit 1
                  fi
                  sleep 1
              fi
          done

          TOPOLOGY=$(curl -s https://{{ .Address }}:{{ .Port }}/topology)
          ETHEREUM_ADDRESS=$(echo ${TOPOLOGY} | jq -r '.self["Ethereum Address"]')

          {{- range nomadService "mev-commit-geth-bootnode1" }}
            {{- if contains "http" .Tags }}
            local/cast send \
                 --rpc-url http://{{ .Address }}:{{ .Port }} \
                 --private-key 0x7c9bf0f015874594d321c1c01ada3166c3509bbd91f76f9e4d7380c2df269c55 ${ETHEREUM_ADDRESS} \
                 --value 100ether
            {{- end }}
          {{- end }}

          if ! curl -s -X POST https://{{ .Address }}:{{ .Port }}/v1/bidder/prepay/1000000000000000000 > /dev/null; then
            echo "Failed to send bidder prepay transaction."
            exit 1
          fi
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
  }
}
