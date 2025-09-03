FROM golang:1.24-alpine AS build

WORKDIR /ws

COPY go.work go.work.sum ./

COPY contracts-abi/go.mod    contracts-abi/go.sum   ./contracts-abi/
COPY p2p/go.mod              p2p/go.sum             ./p2p/
COPY oracle/go.mod           oracle/go.sum          ./oracle/
COPY testing/go.mod          testing/go.sum         ./testing/
COPY tools/go.mod            tools/go.sum           ./tools/
COPY x/go.mod                x/go.sum               ./x/
COPY bridge/standard/go.mod  bridge/standard/go.sum ./bridge/standard/
COPY cl/go.mod               cl/go.sum              ./cl/
COPY infrastructure/tools/keystore-generator/go.mod infrastructure/tools/keystore-generator/go.sum ./infrastructure/tools/keystore-generator/

COPY p2p/integrationtest/provider/entrypoint.sh /scripts/provider-emulator-entrypoint.sh

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go work sync && \
    go mod download all

COPY . .

ARG TARGETS="./oracle/cmd \
             ./p2p/cmd \
             ./bridge/standard/cmd/relayer \
             ./bridge/standard/cmd/emulator \
             ./infrastructure/tools/keystore-generator \
             ./testing/cmd \
             ./tools/preconf-rpc \
             ./tools/beacon-emulator \
             ./tools/dashboard \
             ./tools/bidder-cli \
             ./tools/bls-signer \
             ./tools/bidder-emulator \
             ./tools/relay-emulator \
             ./tools/validators-monitor \
             ./tools/points-service \
             ./p2p/integrationtest/provider"

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    set -e; \
    for path in ${TARGETS}; do \
        # If the last element is literally "cmd", use the directory above it
        bn=$(basename "$path"); \
        if [ "$bn" = "cmd" ]; then \
            name=$(basename "$(dirname "$path")"); \
        else \
            name=$bn; \
        fi; \
        echo "→ building $path as /go/bin/$name"; \
        CGO_ENABLED=0 go build -o "/go/bin/$name" "$path"; \
    done
