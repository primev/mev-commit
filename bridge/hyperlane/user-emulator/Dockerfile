FROM ghcr.io/foundry-rs/foundry:latest

RUN apk update && \
    apk add bash jq python3 py3-pip bc

RUN pip3 install datadog

COPY hyperlane/cli.sh /usr/local/bin/bridge-cli
RUN chmod +x /usr/local/bin/bridge-cli

COPY hyperlane/user-emulator/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
