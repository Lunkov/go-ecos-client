package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
)

func TestMsgTransaction(t *testing.T) {
  w1 := wallets.NewWallet(wallets.TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  w2 := wallets.NewWallet(wallets.TypeWalletHD)
  w2.Create(&map[string]string{"mnemonic": "fall farm prepare palm sign city analyst liquid orange naive hire lawn marble object old cradle exchange visa caught base robot online undo possible"})

  addr1 := w1.GetAddress(hdwallet.ECOS)
  addr2 := w2.GetAddress(hdwallet.ECOS)
  
  assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", addr1)
  assert.Equal(t, "0xfa242EE498857ec3C06a2E5E9e37b090807B467a", addr2)
  
  obj := NewMsgTransaction()
  obj.Init(1, w1, addr2, hdwallet.ECOS, hdwallet.ECOS, 0, 10)
  
  oks := obj.DoSign(w1)
  assert.True(t, oks)

  _, ok := obj.ToJSON()
  assert.True(t, ok)
  /*
   * assert.Equal(t, "{\"version\":\"\",\"id_message\":2907279365,\"id_object\":\"payment\",\"id_action\":1,\"updated_at\":\"2023-05-24T10:17:54.256548986+03:00\",\"pubkey\":null,\"sign\":null,\"address_from\":\"0x5fAD534AadacBe64E43944CAEfAC04B087B75F9D\",\"address_to\":\"0x8E348b74e2f52f1c97ADBf0aA42b4c3FC7961fA6\",\"coin_from\":2147685952,\"coin_to\":2147685952}", bufJson)
  */
  
  buf := obj.Serialize()
  
  obj2 := NewMsgTransaction()
  obj2.Deserialize(buf)

  assert.True(t, obj2.DoVerify()) 
}
