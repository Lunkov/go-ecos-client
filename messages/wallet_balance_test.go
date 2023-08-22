package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
)

func TestMsgWalletBalance(t *testing.T) {
  w1 := wallets.NewWallet(wallets.TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  msg := NewGetBalanceReq()
  err := msg.Init(w1, hdwallet.ECOS)
  assert.Nil(t, err)
  
  buf := msg.Serialize()
  
  msg2 := NewGetBalanceReq()
  msg2.Deserialize(buf)

  vok, _ := msg2.DoVerify()
  assert.True(t, vok) 
}
