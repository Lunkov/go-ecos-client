package client

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
  
  "github.com/Lunkov/go-ecos-client/objects"
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
    assert.Equal(t, uint64(27), tx.GetValueOut())
  }
  
  passport := messages.NewPassportInfo()
  passport.Passport.DisplayName = "User1 Tester" 
  passport.Passport.FirstName = "User1" 
  passport.Passport.LastName = "Tester"
  passport.Passport.Photo = []byte{0x01, 0x02}
  passport.CodePhrase = "password" 
  passport.TX = *tx
  
  tx, ok = client.PassportCommit(w1, hdwallet.ECOS, passport)
  assert.True(t, ok)
  
  if tx != nil {
    assert.Equal(t, objects.Transaction{Id:[]uint8{}, Version:0x0, Timestamp:1686598557, Vin:[]objects.TXInput{objects.TXInput{Txid:[]uint8(nil), Address:"0x25623Fc60db8bDBB984f2Bb1246ef752e9CC5c41", CIDNFT:"", Vout:0x0, Signature:[]uint8{}, PublicKey:[]uint8{}}, objects.TXInput{Txid:[]uint8{}, Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"", Vout:0x1b, Signature:[]uint8(nil), PublicKey:[]uint8(nil)}}, Vout:[]objects.TXOutput{objects.TXOutput{Address:"0x5f7ae710cED588D42E863E9b55C7c51e56869963", CIDNFT:"passport:QmYANktUduDRGvuM6wPqM84gzsZJWjKQYDCrXjB4XLnvKC", Value:0x0}, objects.TXOutput{Address:"0x25623Fc60db8bDBB984f2Bb1246ef752e9CC5c41", CIDNFT:"", Value:0x14}, objects.TXOutput{Address:"0xe414f133160Eced6e00CF686f97c19809803EF04", CIDNFT:"", Value:0x5}, objects.TXOutput{Address:"0xfa242EE498857ec3C06a2E5E9e37b090807B467a", CIDNFT:"", Value:0x2}}}, *tx)
    
    p2, pok := client.PassportGet(tx.Vout[0].CIDNFT, passport.CodePhrase)
    assert.True(t, pok)
    
    if p2 != nil {
      assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", *p2)
    }
  }
  
}
