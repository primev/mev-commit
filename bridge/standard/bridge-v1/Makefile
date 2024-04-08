.PHONY: relayer user_cli
export CGO_ENABLED=0

relayer: bin
	go build -o bin/relayer ./cmd/relayer

linux-amd64-relayer:
	GOOS=linux GOARCH=amd64 go build -o mev-commit-bridge-relayer-linux-amd64 ./cmd/relayer
	tar -czvf mev-commit-bridge-relayer-linux-amd64.tar.gz mev-commit-bridge-relayer-linux-amd64

darwin-arm64-relayer:
	GOOS=darwin GOARCH=arm64 go build -o mev-commit-bridge-relayer-darwin-arm64 ./cmd/relayer
	tar -czvf mev-commit-bridge-relayer-darwin-arm64.tar.gz mev-commit-bridge-relayer-darwin-arm64

user-cli: bin
	go build -o bin/user_cli ./cmd/user_cli

user-emulator: bin
	go build -o bin/user_emulator ./emulator

bin:
	mkdir -p $@

up:
	docker-compose up --build -d

up-agent:
	docker-compose --profile dd_agent up --build -d

down:
	docker-compose --profile dd_agent down
