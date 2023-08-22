module github.com/Lunkov/go-ecos-client/messages

go 1.19

replace go-ecos-client/utils => ./../utils

replace go-ecos-client/objects => ./../objects

require (
	github.com/Lunkov/go-hdwallet v0.0.0-20230525092819-390711df8fa3
	github.com/Lunkov/lib-wallets v0.0.0-20230822101559-a0b70aca6d9f
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/Lunkov/go-btcec v0.0.0-20230525101159-f058a4a0edc0 // indirect
	github.com/Lunkov/lib-cipher v0.0.0-20230822094401-480fc8192b31 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/cpacia/bchutil v0.0.0-20181003130114-b126f6a35b6c // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/go-ethereum v1.12.2 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	go-ecos-client/objects v0.0.0-00010101000000-000000000000 // indirect
	go-ecos-client/utils v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
