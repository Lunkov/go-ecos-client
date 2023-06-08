package objects

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/lib-wallets"
)

func TestTx(t *testing.T) {
  w1 := wallets.NewWallet(wallets.TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  obj := NewTX()
  oks := obj.DoSign(w1)
  assert.True(t, oks)

  buf, ok2 := obj.Serialize()
  assert.True(t, ok2)
  
  obj2 := NewTX()
  obj2.Deserialize(buf)

  assert.True(t, obj2.DoVerify()) 
}
