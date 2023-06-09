package client

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
  
  "github.com/Lunkov/go-ecos-client/messages"
)

func TestPassport(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  //flag.Set("v", "9")
  //flag.Parse()

  client := NewClientECOS([]string{"http://127.0.0.1:8085/"}, 3)

  w1 := wallets.NewWallet(wallets.TypeWalletHD)
  w1.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})


  tx, ok := client.PassportNew(w1, hdwallet.ECOS)
  assert.True(t, ok)
  
  if tx != nil {
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", *tx)
  }
  
  passport := messages.NewPassportInfo()
  passport.Passport.DisplayName = "User1 Tester" 
  passport.Passport.FirstName = "User1" 
  passport.Passport.LastName = "Tester"
  
  passport.CodePhrase = "password" 
  
  tx, ok = client.PassportCommit(w1, hdwallet.ECOS, passport)
  assert.True(t, ok)
  
  if tx != nil {
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", *tx)
    
    p2, pok := client.PassportGet(tx.Vout[0].CIDNFT, passport.CodePhrase)
    assert.True(t, pok)
    
    if p2 != nil {
      assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", *p2)
    }
  }
  
}
