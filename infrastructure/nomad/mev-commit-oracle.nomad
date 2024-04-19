job "mev-commit-oracle" {
  datacenters = ["dc1"]

  meta {
    POSTGRESQL_MAIN_VERSION = "15"
    DB_NAME="mev_oracle"
    DB_USER="mev_oracle"
    MEV_COMMIT_ORACLE_VERSION = "v0.3.0-rc1"
    KEYSTORE_GENERATOR_VERSION_ID = "RYncrfSOWyiyojCLlse4f4kSiotrdhgM"
  }

  group "mev-commit-oracle-group" {
    count = 1

    network {
      mode = "bridge"

      dns {
        servers = ["127.0.0.53", "1.1.1.1", "8.8.8.8", "8.8.4.4"]
      }

      port "db" {
        static = 5432
        to = 5432
      }
      port "http" {
        to = 8080
      }
    }

    service {
      name = "mev-commit-oracle-db"
      port = "db"
      tags = ["db"]
      provider = "nomad"
    }

    service {
      name = "mev-commit-oracle"
      port = "http"
      tags = ["http"]
      provider = "nomad"
    }

    task "db" {
      driver = "exec"

      lifecycle {
        hook    = "prestart"
        sidecar = true
      }

      template {
        data = <<-EOH
          PATH = "/usr/lib/postgresql/{{env "NOMAD_META_POSTGRESQL_MAIN_VERSION"}}/bin:{{env "PATH"}}"
          PG_DATA = "/local/pgdata-{{env "NOMAD_ALLOC_INDEX"}}"
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

          if [ -d "${PG_DATA}" ]; then
              "Initialized and configured database found"
              postgres -D ${PG_DATA}
              exit $?
          fi

          mkdir -p /var/run/postgresql > /dev/null 2>&1
          pg_ctl initdb --silent --pgdata=${PG_DATA}
          postgres -D ${PG_DATA} &
          pid=$!
          until pg_isready --quiet --username=$USER --dbname=postgres; do
            echo "Waiting for PostgreSQL to start..."
            sleep 1
          done

          PASSWORD="{{ with secret "secret/data/mev-commit" }}{{ .Data.data.oracle_db_password }}{{ end }}"
          createuser --username=$USER --createdb ${NOMAD_META_DB_USER}
          createdb --username=${NOMAD_META_DB_USER} ${NOMAD_META_DB_NAME}
          psql --quiet \
               --username=${NOMAD_META_DB_USERNAME} \
               --dbname=${NOMAD_META_DB_NAME} \
               --command="ALTER USER ${NOMAD_META_DB_USER} WITH PASSWORD '${PASSWORD}'; \
                          GRANT ALL PRIVILEGES ON DATABASE ${NOMAD_META_DB_NAME} TO ${NOMAD_META_DB_USER};"
          echo "Database initialized and configured successfully"

          wait $pid
        EOH
        destination = "local/run.sh"
        perms = "0755"
      }

      config {
        command = "bash"
        args = ["-c", "local/run.sh"]
      }
    }

    task "oracle" {
      driver = "exec"

      artifact {
        source = "https://github.com/primevprotocol/mev-commit-oracle/releases/download/${NOMAD_META_MEV_COMMIT_ORACLE_VERSION}/mev-commit-oracle_Linux_x86_64.tar.gz"
      }

      artifact {
        source = "https://primev-infrastructure-artifacts.s3.us-west-2.amazonaws.com/keystore-generator?versionId=${NOMAD_META_KEYSTORE_GENERATOR_VERSION_ID}"
      }

      template {
        data = <<-EOH
          KEYSTOREGEN_LOG_FMT="json"
          KEYSTOREGEN_LOG_TAGS="service:mev-commit-oracle"
          MEV_ORACLE_LOG_FMT="json"
          MEV_ORACLE_LOG_TAGS="service:mev-commit-oracle"
          MEV_ORACLE_KEYSTORE_PATH = "/local/data-{{env "NOMAD_ALLOC_INDEX"}}/keystore"
          MEV_ORACLE_KEYSTORE_PASSWORD = "{{ with secret "secret/data/mev-commit" }}{{ .Data.data.oracle_keystore_password }}{{ end }}"
          MEV_ORACLE_L1_RPC_URL = "https://rpc.sepolia.org"
      {{- range nomadService "mev-commit-oracle" }}
          {{- if contains "http" .Tags }}
          MEV_ORACLE_HTTP_PORT = "{{ .Port }}"
          {{- end }}
      {{- end }}
      {{- range nomadService "mev-commit-geth-bootnode1" }}
          {{- if contains "http" .Tags }}
          MEV_ORACLE_SETTLEMENT_RPC_URL = "http://{{ .Address }}:{{ .Port }}"
          {{- end }}
      {{- end }}
      {{- range nomadService "mev-commit-oracle-db" }}
          {{- if contains "db" .Tags }}
          MEV_ORACLE_PG_HOST = "localhost"
          MEV_ORACLE_PG_PORT = "{{ .Port }}"
          MEV_ORACLE_PG_USER = "{{env "NOMAD_META_DB_USER"}}"
          MEV_ORACLE_PG_PASSWORD = "{{ with secret "secret/data/mev-commit" }}{{ .Data.data.oracle_db_password }}{{ end }}"
          MEV_ORACLE_PG_DBNAME = "{{env "NOMAD_META_DB_NAME"}}"
          {{- end }}
      {{- end }}
          MEV_ORACLE_LOG_LEVEL = "info"
          MEV_ORACLE_LAGGERD_MODE = 64
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

          if [ ! -d "${MEV_ORACLE_KEYSTORE_PATH}" ]; then
            mkdir -p "${MEV_ORACLE_KEYSTORE_PATH}" > /dev/null 2>&1
            chmod +x local/keystore-generator
            local/keystore-generator generate \
              --keystore-dir "${MEV_ORACLE_KEYSTORE_PATH}" \
              --passphrase "${MEV_ORACLE_KEYSTORE_PASSWORD}"
          fi

          mkdir -p $(dirname "${MEV_ORACLE_PRIV_KEY_FILE}") > /dev/null 2>&1
          chmod +x local/mev-commit-oracle
          local/mev-commit-oracle start
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
