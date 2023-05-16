package objects

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/lib-messages"
  "github.com/Lunkov/lib-wallets"
)

func TestRule(t *testing.T) {
	client := NewClient([]string{"http://127.0.0.1:8084/"}, 3)
  
  w := &wallets.WalletHD{}
  w.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})
  balance, ok := client.GetBalance(w)
	assert.True(t, ok)
  assert.Equal(t, &messages.Balance{Address: "0x5f7ae710cED588D42E863E9b55C7c51e56869963"}, balance)
  
  
  
}
