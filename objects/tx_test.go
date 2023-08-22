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
  serr := obj.DoSign(w1)
  assert.Nil(t, serr)

  buf, errs := obj.Serialize()
  assert.Nil(t, errs)
  
  obj2 := NewTX()
  obj2.Deserialize(buf)

  vok, verr := obj2.DoVerify()
  assert.True(t, vok)
  assert.Nil(t, verr) 
}
