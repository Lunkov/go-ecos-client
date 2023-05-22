package client

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/lib-wallets"
  
  // "github.com/Lunkov/go-ecos-client/messages"
)

func TestCoins(t *testing.T) {
	client := NewClientECOS([]string{"http://127.0.0.1:8085/"}, 3)
  
  w := &wallets.WalletHD{}
  w.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})
  balance, ok := client.GetBalance(w)
  assert.True(t, ok)
  
  if balance != nil {
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", balance.Address)
    assert.Equal(t, "", balance.Coin)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.Balance)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.UnconfirmedBalance)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.TotalReceived)
    assert.Equal(t, uint64(0x00), balance.TotalSent)
  }
  
  balance, ok = client.NewTransaction(w, "0xe414f133160Eced6e00CF686f97c19809803EF04", "ECOS", 1000, 10)
	assert.True(t, ok)
  if balance != nil {
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", balance.Address)
    assert.Equal(t, "", balance.Coin)
    assert.Equal(t, 0x2386f26fc10000, balance.Balance)
    assert.Equal(t, 0x2386f26fc10000, balance.UnconfirmedBalance)
    assert.Equal(t, 0x2386f26fc10000, balance.TotalReceived)
    assert.Equal(t, 0x00, balance.TotalSent)
  }
}
