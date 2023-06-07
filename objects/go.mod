module github.com/Lunkov/go-ecos-client/objects

go 1.19

replace github.com/Lunkov/go-ecos-client/utils => ./../utils

require (
	github.com/Lunkov/go-ecos-client/utils v0.0.0-20230606195050-4a0830ead330
	github.com/Lunkov/lib-cipher v0.0.0-20230420102046-39f2f16e9d29
	github.com/golang/glog v1.1.1
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/Lunkov/go-btcec v0.0.0-20230525101159-f058a4a0edc0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/ethereum/go-ethereum v1.12.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
