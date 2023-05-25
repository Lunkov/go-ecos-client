module github.com/Lunkov/go-ecos-client

go 1.19

replace github.com/Lunkov/go-ecos-client/messages => ./messages

replace github.com/Lunkov/go-ecos-client/utils => ./utils

replace github.com/Lunkov/go-ecos-client/objects => ./objects

require (
	github.com/Lunkov/go-ecos-client/messages v0.0.0-20230525065821-e4e16eb08efb
	github.com/Lunkov/go-hdwallet v0.0.0-20230402114829-9836cf5dfed5
	github.com/Lunkov/lib-wallets v0.0.0-20230525055837-13776bf398ad
	github.com/golang/glog v1.1.1
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/Lunkov/go-ecos-client/utils v0.0.0-20230525065821-e4e16eb08efb // indirect
	github.com/Lunkov/lib-cipher v0.0.0-20230324195628-77d817f26180 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/cpacia/bchutil v0.0.0-20181003130114-b126f6a35b6c // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/go-ethereum v1.11.5 // indirect
	github.com/holiman/uint256 v1.2.0 // indirect
	github.com/itchyny/base58-go v0.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
