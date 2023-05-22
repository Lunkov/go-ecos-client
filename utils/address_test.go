package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  // "crypto/sha512"
  // "crypto/ecdsa"
  "github.com/Lunkov/lib-wallets"
)

func TestAddress(t *testing.T) {
  
  w := &wallets.WalletHD{}
  w.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", w.GetAddress("ECOS"))
  assert.Equal(t, "", PubkeyToAddress(w.Master.PublicECDSA).Hex())
}

