package clientecos

import (
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-messages"
  "github.com/Lunkov/lib-wallets"
)

func (c *Client) GetBalance(w *wallets.WalletHD) (*messages.Balance, bool) {
  c.selectServer()
  msg := messages.NewReqGetBalance()
  msg.Address = w.GetAddress("ECOS")
  
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

func (c *Client) NewTransaction(w *wallets.WalletHD, addressTo string, coin string, value uint64, maxCost uint64) (*messages.Balance, bool) {
  c.selectServer()
  msg := messages.NewTokenTransaction()
  pkBuf, okpk := PublicKeyToBytes(w.Master.PrivateECDSA)
  if !okpk {
    glog.Errorf("ERR: NewTransaction.PublicKeyToBytes")
    return nil, false
  }
  msg.Init(w.GetAddress("ECOS"), addressTo, coin, value, maxCost, pkBuf)
  
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
