package client

import (
  "github.com/golang/glog"
  
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
  "github.com/Lunkov/go-ecos-client/messages"
  "github.com/Lunkov/go-ecos-client/utils"
)

func (c *ClientECOS) GetBalance(w wallets.IWallet) (*messages.Balance, bool) {
  c.selectServer()
  msg := messages.NewReqGetBalance()
  msg.Address = w.GetAddress(hdwallet.ECOS)
  
  answer, ok := c.httpRequest("/wallet/balance", string(msg.Serialize()))
  if !ok {
    return nil, false
  }
  result := messages.NewBalance()
  if !result.Deserialize(answer) {
    glog.Errorf("ERR: GetBalance.Deserialize")
    return nil, false
  }
  return result, true
}

func (c *ClientECOS) NewTransaction(w wallets.IWallet, addressTo string, coin uint32, value uint64, maxCost uint64) (*messages.Balance, bool) {
  c.selectServer()
  msg := messages.NewTokenTransaction()
  pkBuf, okpk := utils.ECDSAPublicKeySerialize(w.GetECDSAPublicKey())
  if !okpk {
    glog.Errorf("ERR: NewTransaction.PublicKeyToBytes")
    return nil, false
  }
  msg.Init(w.GetAddress(hdwallet.ECOS), addressTo, coin, value, maxCost, pkBuf)
  
  answer, ok := c.httpRequest("/new/transaction", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  result := messages.NewBalance()
  if !result.Deserialize(answer) {
    glog.Errorf("ERR: NewTransaction.Deserialize")
    return nil, false
  }
  return result, true
}
