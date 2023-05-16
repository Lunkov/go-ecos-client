package utils

import (
  "crypto/sha256"
  "github.com/itchyny/base58-go"
  "golang.org/x/crypto/ripemd160"
  "github.com/golang/glog"
)

// ReverseBytes reverses a byte array
func ReverseBytes(data *[]byte) {
  for i, j := 0, len(*data)-1; i < j; i, j = i+1, j-1 {
    (*data)[i], (*data)[j] = (*data)[j], (*data)[i]
  }
}

func UInt64ToBytes(v uint64) (b []byte) {
  l := 8
  b = make([]byte, l)

  for i := 0; i < l; i++ {
    f := 8 * i
    b[i] = byte(v >> f) 
  }
  return b
}

func UInt32ToBytes(v uint32) (b []byte) {
  l := 4
  b = make([]byte, l)

  for i := 0; i < l; i++ {
    f := 8 * i
    b[i] = byte(v >> f)
  }
  return b
}

func Int64ToBytes(v int64) (b []byte) {
  l := 8
  b = make([]byte, l)

  for i := 0; i < l; i++ {
    f := 8 * i
    b[i] = byte(v >> f)
  }
  return b
}

// GetAddress returns wallet address
func GetAddressPublicKey(publicKey []byte) []byte {
  pubKeyHash := HashPubKey(publicKey)

  versionedPayload := append([]byte{version}, pubKeyHash...)
  checksum := checksum(versionedPayload)

  fullPayload := append(versionedPayload, checksum...)
  
  encoding := base58.BitcoinEncoding // or RippleEncoding or BitcoinEncoding
  address, _ := encoding.Encode(fullPayload)

  return address
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
  publicSHA256 := sha256.Sum256(pubKey)

  RIPEMD160Hasher := ripemd160.New()
  _, err := RIPEMD160Hasher.Write(publicSHA256[:])
  if err != nil {
    glog.Errorf("ERR: RIPEMD160Hasher: %v", err)
  }
  publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

  return publicRIPEMD160
}

// Checksum generates a checksum for a public key
const addressChecksumLen = 4
//
func checksum(payload []byte) []byte {
  firstSHA := sha256.Sum256(payload)
  secondSHA := sha256.Sum256(firstSHA[:])

  return secondSHA[:addressChecksumLen]
}

