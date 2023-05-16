package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "crypto/sha512"
  "crypto/ecdsa"
  "github.com/Lunkov/lib-wallets"
)

func TestRSA(t *testing.T) {
  
  w := &wallets.WalletHD{}
  w.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  pkBuf, okpk := RSASerialize(w.Master.PublicECDSA)

	assert.True(t, okpk)
  assert.Equal(t, []byte{0x22, 0xff, 0x81, 0x3, 0x1, 0x1, 0x8, 0x45, 0x43, 0x44, 0x53, 0x41, 0x42, 0x75, 0x66, 0x1, 0xff, 0x82, 0x0, 0x1, 0x2, 0x1, 0x1, 0x58, 0x1, 0xa, 0x0, 0x1, 0x1, 0x59, 0x1, 0xa, 0x0, 0x0, 0x0, 0x47, 0xff, 0x82, 0x1, 0x20, 0x23, 0xf, 0x74, 0x93, 0x9e, 0x44, 0x5, 0x58, 0xc0, 0xf5, 0xf8, 0xc7, 0x37, 0xeb, 0x8, 0x7b, 0xd6, 0x6, 0x25, 0x60, 0xb0, 0x62, 0xc6, 0x9a, 0x18, 0x68, 0xef, 0x18, 0x5, 0x68, 0xe1, 0x6b, 0x1, 0x20, 0xc, 0x74, 0x19, 0x99, 0xe7, 0x3a, 0x8a, 0xd9, 0xdf, 0x6c, 0xd3, 0xb3, 0x56, 0x11, 0x4b, 0x4c, 0x1, 0x54, 0x7e, 0xf7, 0x85, 0x84, 0x71, 0xc, 0x32, 0xb2, 0xe6, 0xc7, 0x6a, 0x9, 0x49, 0x53, 0x0}, pkBuf)
  
  message := []byte("Hello world")
  signature, sok := RSASign(w.Master.PrivateECDSA, message)
  assert.True(t, sok)
  
  hashed := sha512.Sum512(message)
  vok1 := ecdsa.VerifyASN1(w.Master.PublicECDSA, hashed[:], signature)
  assert.True(t, vok1)
}
