module github.com/iotaledger/wasp/tools/gascalibration

go 1.21

replace (
	github.com/ethereum/go-ethereum => github.com/iotaledger/go-ethereum v1.12.0-wasp
	github.com/iotaledger/wasp => ../../
	github.com/iotaledger/wasp/tools/wasp-cli => ../wasp-cli/
	go.dedis.ch/kyber/v3 => github.com/kape1395/kyber/v3 v3.0.14-0.20230124095845-ec682ff08c93 // branch: dkg-2suites
)

require (
	github.com/iotaledger/wasp/tools/wasp-cli v1.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.7.0
	gonum.org/v1/plot v0.14.0
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	git.sr.ht/~sbinet/gg v0.5.0 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.1 // indirect
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cockroachdb/errors v1.11.0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v0.0.0-20230412222916-60cfeb46143b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set/v2 v2.3.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/ethereum/go-ethereum v1.12.2 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/getsentry/sentry-go v0.23.0 // indirect
	github.com/go-fonts/liberation v0.3.1 // indirect
	github.com/go-latex/latex v0.0.0-20230307184459-12ec69307ad9 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-pdf/fpdf v0.8.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/huin/goupnp v1.2.0 // indirect
	github.com/iancoleman/orderedmap v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/iotaledger/hive.go/constraints v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/crypto v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/ds v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/kvstore v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/lo v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/logger v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/runtime v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/serializer/v2 v2.0.0-rc.1.0.20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/hive.go/stringify v0.0.0-20230629181801-64c530ff9d15 // indirect
	github.com/iotaledger/iota.go v1.0.0 // indirect
	github.com/iotaledger/iota.go/v3 v3.0.0-rc.3 // indirect
	github.com/iotaledger/wasp v1.0.0-00010101000000-000000000000 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/oasisprotocol/ed25519 v0.0.0-20210505154701-76d8c688d86e // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.6 // indirect
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	github.com/petermattis/goid v0.0.0-20230808133559-b036b712a89b // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/samber/lo v1.38.1 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/status-im/keycard-go v0.2.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.dedis.ch/fixbuf v1.0.3 // indirect
	go.dedis.ch/kyber/v3 v3.1.0 // indirect
	go.uber.org/goleak v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20231005195138-3e424a577f31 // indirect
	golang.org/x/image v0.11.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
