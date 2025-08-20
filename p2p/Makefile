LDFLAGS ?= -s -w \
-X github.com/primev/mev-commit/p2p.version=$(shell git describe --tags || git rev-parse --short HEAD)

.PHONY: build
build: export CGO_ENABLED=0
build: bin
	go build -ldflags '$(LDFLAGS)' -o bin/mev-commit ./cmd

bin:
	mkdir $@

bufgen:
	cd rpc && buf generate -o ../
	cd messages && buf generate -o ../

bidder:
	export MEV_COMMIT_PRIVKEY_FILE=$(realpath cmd/defaultkey)
	go run cmd/main.go \
		--peer-type=bidder \
		--bidder-registry-contract="0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9" \
		--block-tracker-contract="0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512" \
		--oracle-contract="0x0DCd1Bf9A1b36cE34237eEaFef220932846BCD82" \
		--preconf-contract="0x8A791620dd6260079BF849Dc5567aDC3F2FdC318" \
		--provider-registry-contract="0xa513E6E4b8f2a923D98304ec87F64353C4D5C853" \
		--gas-tip-cap=1000000000 \
		--gas-fee-cap=3000000000
