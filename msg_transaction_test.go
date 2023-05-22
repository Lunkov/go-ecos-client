package client

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/lib-wallets"
)

func TestMsgTransaction(t *testing.T) {
  w1 := &wallets.WalletHD{}
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  w2 := &wallets.WalletHD{}
  w2.Create(&map[string]string{"mnemonic": "elegant chase oil pigeon ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})

  obj := NewMsgTransaction()
  obj.Init(1, w1.GetAddress("ecos"), w2.GetAddress("ecos"), "ecos", "ecos", 0, 10)
  
  assert.Equal(t, []byte{0x52, 0xda, 0x34, 0x3, 0x4e, 0x18, 0xbb, 0xa, 0x3, 0x7c, 0x1d, 0x79, 0xe9, 0x36, 0x55, 0xde, 0xc3, 0xd0, 0x38, 0x9f, 0xb6, 0xe, 0x6f, 0x9c, 0x39, 0x15, 0xb4, 0x11, 0x28, 0xc5, 0xcf, 0x7f, 0xd5, 0x8c, 0x9, 0xea, 0xff, 0x9b, 0x63, 0xdd, 0x6a, 0x72, 0xf1, 0xea, 0x91, 0x1a, 0x44, 0x15, 0xb8, 0x7e, 0x4f, 0x52, 0x28, 0x0, 0xa9, 0x72, 0xa0, 0xfc, 0xc6, 0x7f, 0xae, 0xe, 0x7b, 0x75}, obj.Hash())
  assert.Equal(t, uint32(0xda46c90), obj.msgId())
  
  buf := obj.Serialize()
  
  obj2 := NewMsgTransaction()
  obj2.Deserialize(buf)

  bufJson, ok := obj.ToJSON()
  assert.True(t, ok)
  assert.Equal(t, "{\"version\":\"\",\"id_message\":2483562631,\"id_object\":\"token\",\"id_action\":1,\"updated_at\":\"2023-05-20T13:50:44.183894305+03:00\",\"pubkey\":null,\"sign\":null,\"address_from\":\"\",\"address_to\":\"\",\"coin_from\":\"\",\"coin_to\":\"\"}", bufJson)

  oks := obj.DoSign(w1.Master.PrivateECDSA)
  assert.True(t, oks)
}
