package utils

import (
  "hash"
  "golang.org/x/crypto/sha3"
  "crypto/elliptic"
  "crypto/ecdsa"
  "encoding/hex"
)

// https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go

/////////// Address
// Lengths of hashes and addresses in bytes.
const (
  // HashLength is the expected length of the hash
  HashLength = 32
  // AddressLength is the expected length of the address
  AddressLength = 20
)

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
  var a Address
  a.SetBytes(b)
  return a
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a), b will be cropped from the left.
func (a *Address) SetBytes(b []byte) {
  if len(b) > len(a) {
    b = b[len(b)-AddressLength:]
  }
  copy(a[AddressLength-len(b):], b)
}


// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	return string(a.checksumHex())
}

func (a *Address) checksumHex() []byte {
	buf := a.hex()

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:]
}

func (a Address) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
  hash.Hash
  Read([]byte) (int, error)
}

// NewKeccakState creates a new KeccakState
func newKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
  b := make([]byte, 32)
  d := newKeccakState()
  for _, b := range data {
    d.Write(b)
  }
  d.Read(b)
  return b
}

func fromECDSAPub(pub *ecdsa.PublicKey) []byte {
  if pub == nil || pub.X == nil || pub.Y == nil {
    return nil
  }
  return elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)
}

func PubkeyToAddress(p ecdsa.PublicKey) Address {
  pubBytes := fromECDSAPub(&p)
  return BytesToAddress(keccak256(pubBytes[1:])[12:])
}
