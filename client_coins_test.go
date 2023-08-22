package client

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
  
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
  
  "go-ecos-client/objects"
)

func TestCoins(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  //flag.Set("v", "9")
  //flag.Parse()

  client := NewClientECOS([]string{"http://127.0.0.1:8085/"}, 3)
  clientLong := NewClientECOS([]string{
                                  "http://127.0.0.1:8081/",
                                  "http://127.0.0.1:8082/",
                                  "http://127.0.0.1:8083/",
                                  "http://127.0.0.1:8084/",
                                  "http://127.0.0.1:8088/",
                                  "http://127.0.0.1:8086/",
                                  "http://127.0.0.1:8087/",
                                  "http://127.0.0.1:8085/",
                                 }, 10)
  
  w := wallets.NewWallet(wallets.TypeWalletHD)
  w.Create(&map[string]string{"mnemonic": "chase oil pigeon elegant ketchup whip frozen beauty unknown brass amount slender pony pottery attitude flavor rifle primary beach sign glue oven crazy lottery"})
  
  w2 := wallets.NewWallet(wallets.TypeWalletHD)
  w2.Create(&map[string]string{"mnemonic": "super inflict toward range sting apart rebel lady detail exotic remove school love buddy guitar hen loan approve neither giant pig divide glimpse loud"})
  
  balance, err := clientLong.GetBalance(w)
  assert.Nil(t, err)

  if balance != nil {
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", balance.Address)
    assert.Equal(t, hdwallet.ECOS,            balance.Coin)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.Balance)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.UnconfirmedBalance)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.TotalReceived)
    assert.Equal(t, uint64(0x00), balance.TotalSent)
  }

  balance, err = client.GetBalance(w)
  assert.Nil(t, err)
  
  if balance != nil {
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", balance.Address)
    assert.Equal(t, hdwallet.ECOS,            balance.Coin)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.Balance)
    assert.Equal(t, uint64(0), balance.UnconfirmedBalance)
    assert.Equal(t, uint64(0x2386f26fc10000), balance.TotalReceived)
    assert.Equal(t, uint64(0x00), balance.TotalSent)
  }
  
  // NewTransaction(w wallets.IWallet, addressTo string, coin uint32, value uint64)
  newTx, errTx := client.TransactionNew(w, "0xe414f133160Eced6e00CF686f97c19809803EF04", hdwallet.ECOS, 1000)
	assert.Nil(t, errTx)
  if newTx != nil {
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", newTx.Vin)
    assert.Equal(t, []objects.TXOutput([]objects.TXOutput{objects.TXOutput{Address:"0xe414f133160Eced6e00CF686f97c19809803EF04", CIDNFT:"", Value:0x3e8},
                                      objects.TXOutput{Address:"0xe414f133160Eced6e00CF686f97c19809803EF04", CIDNFT:"", Value:0x5},
                                      objects.TXOutput{Address:"0xfa242EE498857ec3C06a2E5E9e37b090807B467a", CIDNFT:"", Value:0x2}}),
                                      newTx.Vout)
    //assert.Equal(t, hdwallet.ECOS,    newTx.CoinFrom)
    //assert.Equal(t, hdwallet.ECOS,    newTx.CoinTo)
    //assert.Equal(t, 0x2386f26fc10000, newTx.ValueFrom)
    //assert.Equal(t, 0x2386f26fc10000, newTx.ValueTo)
    //assert.Equal(t, 0x2386f26fc10000, newTx.IdStatus)
  }
  
  newTx, errTx = client.TransactionCommit(w, newTx)
	assert.Nil(t, errTx)
  if newTx != nil {
    /*
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", newTx.AddressFrom)
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", newTx.AddressTo)
    assert.Equal(t, hdwallet.ECOS,    newTx.CoinFrom)
    assert.Equal(t, hdwallet.ECOS,    newTx.CoinTo)
    assert.Equal(t, 0x2386f26fc10000, newTx.ValueFrom)
    assert.Equal(t, 0x2386f26fc10000, newTx.ValueTo)
    assert.Equal(t, 0x2386f26fc10000, newTx.IdStatus)
    * */
  }

  if newTx != nil {
    statTx, errStatTx := client.TransactionStatus(w, newTx.Id)
    assert.Nil(t, errStatTx)
    assert.Equal(t, []byte{0x00}, statTx.IdTx)
    /*assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", newTx.AddressFrom)
    assert.Equal(t, "0x5f7ae710cED588D42E863E9b55C7c51e56869963", newTx.AddressTo)
    assert.Equal(t, hdwallet.ECOS,    newTx.CoinFrom)
    assert.Equal(t, hdwallet.ECOS,    newTx.CoinTo)
    assert.Equal(t, 0x2386f26fc10000, newTx.ValueFrom)
    assert.Equal(t, 0x2386f26fc10000, newTx.ValueTo)
    assert.Equal(t, 0x2386f26fc10000, newTx.IdStatus)*/
  }

  balance, err = client.GetBalance(w2)
  assert.Nil(t, err)
  
  if balance != nil {
    assert.Equal(t, "0xe414f133160Eced6e00CF686f97c19809803EF04", balance.Address)
    assert.Equal(t, hdwallet.ECOS,            balance.Coin)
    assert.Equal(t, uint64(1000000000), balance.Balance)
    assert.Equal(t, uint64(0), balance.UnconfirmedBalance)
    assert.Equal(t, uint64(10000000), balance.TotalReceived)
    assert.Equal(t, uint64(10000), balance.TotalSent)
  }
}
