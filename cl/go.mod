module github.com/primev/mev-commit/cl

go 1.23.0

toolchain go1.24.0

require (
	github.com/go-redis/redismock/v9 v9.2.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang/mock v1.6.0
	github.com/lib/pq v1.10.9
	github.com/primev/mev-commit/p2p v0.0.1
	github.com/redis/go-redis/v9 v9.6.1
	github.com/urfave/cli/v2 v2.27.5
	golang.org/x/sync v0.11.0
	golang.org/x/tools v0.29.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/bits-and-blooms/bitset v1.20.0 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240616162244-4768e80dfb9a // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.2 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/consensys/bavard v0.1.27 // indirect
	github.com/consensys/gnark-crypto v0.16.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/crate-crypto/go-eth-kzg v1.3.0 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240724233137-53bbb0ceb27a // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/ethereum/c-kzg-4844/v2 v2.1.0 // indirect
	github.com/ethereum/go-verkle v0.2.2 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/gomega v1.30.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/supranational/blst v0.3.14 // indirect
	github.com/tklauser/go-sysconf v0.3.13 // indirect
	github.com/tklauser/numcpus v0.7.0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/ethereum/go-ethereum v1.15.11
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/heyvito/go-leader v0.1.0
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/primev/mev-commit/x v0.0.1
	github.com/prometheus/client_golang v1.19.1
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/stretchr/testify v1.10.0
	github.com/vmihailenco/msgpack/v5 v5.4.1
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/primev/mev-commit/p2p => ../p2p
	github.com/primev/mev-commit/x => ../x
)
