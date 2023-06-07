package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
)

func TestMsgWalletTx(t *testing.T) {
  w1 := wallets.NewWallet(wallets.TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  msg := NewWalletTxBalance(w1.GetAddress(hdwallet.ECOS))
  buf, ok := msg.Serialize()
  assert.True(t, ok) 
  
  msg2 := NewWalletTxEmpty()
  
  assert.True(t, msg2.Deserialize(buf)) 
}
