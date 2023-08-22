package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/lib-wallets"
)

func TestMsgTransactionStatus(t *testing.T) {
  w1 := wallets.NewWallet(wallets.TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  obj := NewMsgTransactionStatus([]byte{0x12, 0x34, 0x56, 0x78})
  errs := obj.DoSign(w1)
  assert.Nil(t, errs)

  buf := obj.Serialize()
  
  obj2 := NewMsgTransactionStatus([]byte{0x00})
  obj2.Deserialize(buf)

  assert.Equal(t, obj, obj2)
  
  vok, _ := obj2.DoVerify()
  assert.True(t, vok) 
}
