module github.com/Lunkov/go-ecos-client/objects

go 1.19

replace github.com/Lunkov/go-ecos-client/utils => ./../utils

require (
	github.com/Lunkov/go-ecos-client/utils v0.0.0-20230624064231-4cb5b646b53a
	github.com/Lunkov/go-hdwallet v0.0.0-20230525092819-390711df8fa3
	github.com/Lunkov/lib-cipher v0.0.0-20230324195628-77d817f26180
	github.com/Lunkov/lib-wallets v0.0.0-20230606135804-f0455259642c
	github.com/golang/glog v1.1.1
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/Lunkov/go-btcec v0.0.0-20230525101159-f058a4a0edc0 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/cpacia/bchutil v0.0.0-20181003130114-b126f6a35b6c // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/go-ethereum v1.11.6 // indirect
	github.com/holiman/uint256 v1.2.2-0.20230321075855-87b91420868c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
