module github.com/primev/mev-commit/testing

go 1.22

require (
	github.com/armon/go-radix v1.0.0
	github.com/ethereum/go-ethereum v1.13.14
	github.com/primev/mev-commit/contracts-abi v0.0.1
	github.com/primev/mev-commit/p2p v0.0.1
	github.com/primev/mev-commit/x v0.0.1
	github.com/urfave/cli/v2 v2.27.1
	golang.org/x/sync v0.7.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.2
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.32.0-20240221180331-f05a6f4403ce.1 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.10.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.12.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/crate-crypto/go-kzg-4844 v0.7.0 // indirect
	github.com/deckarep/golang-set/v2 v2.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/ethereum/c-kzg-4844 v0.4.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/supranational/blst v0.3.11 // indirect
	github.com/tklauser/go-sysconf v0.3.13 // indirect
	github.com/tklauser/numcpus v0.7.0 // indirect
	github.com/xrash/smetrics v0.0.0-20231213231151-1d8dd44e695e // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240515191416-fc5f0ca64291 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)

replace github.com/primev/mev-commit/p2p => ../p2p

replace github.com/primev/mev-commit/contracts-abi => ../contracts-abi

replace github.com/primev/mev-commit/x => ../x