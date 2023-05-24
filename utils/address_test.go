package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
)

func TestAddress(t *testing.T) {
  
  w := wallets.NewWallet(wallets.TypeWalletHD)
  w.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", w.GetAddress(hdwallet.ECOS))
  // TODO
  //assert.Equal(t, "", PubkeyToAddress(w.GetECDSAPublicKey()).Hex())
}

