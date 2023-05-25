package client

import (
  "github.com/golang/glog"
  
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"

  "github.com/Lunkov/go-ecos-client/messages"
)

func (c *ClientECOS) GetBalance(w wallets.IWallet) (*messages.Balance, bool) {
  c.selectServer()
  msg := messages.NewGetBalanceReq()
  oks := msg.Init(w, hdwallet.ECOS)
  if !oks {
    return nil, false
  }
  
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

func (c *ClientECOS) TransactionNew(w wallets.IWallet, addressTo string, coin uint32, value uint64) (*messages.MsgTransaction, bool) {
  c.selectServer()
  msg := messages.NewMsgTransaction()
  msg.Init(messages.StatusTxNew, w, addressTo, coin, coin, 0, value)
  
  if !msg.DoSign(w) {
    return nil, false
  }
  
  answer, ok := c.httpRequest("/transaction/new", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  if !msg.Deserialize(answer) {
    glog.Errorf("ERR: NewTransaction.Deserialize")
    return nil, false
  }
  return msg, true
}

func (c *ClientECOS) TransactionStatus(w wallets.IWallet, IdTx uint32) (*messages.MsgTransaction, bool) {
  c.selectServer()
  msg := messages.NewMsgTransaction()
  msg.Init(messages.StatusTxNew, w, "", 0, 0, 0, 0)
  msg.IdTx = IdTx
  
  if !msg.DoSign(w) {
    return nil, false
  }
  
  answer, ok := c.httpRequest("/transaction/status", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  if !msg.Deserialize(answer) {
    glog.Errorf("ERR: NewTransaction.Deserialize")
    return nil, false
  }
  return msg, true
}

func (c *ClientECOS) TransactionCommit(w wallets.IWallet, tx *messages.MsgTransaction) (*messages.MsgTransaction, bool) {
  if tx == nil {
    return nil, false
  }

  c.selectServer()
  
  msg := tx
  
  if !msg.DoSign(w) {
    return nil, false
  }
  
  answer, ok := c.httpRequest("/transaction/new", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  if !msg.Deserialize(answer) {
    glog.Errorf("ERR: NewTransaction.Deserialize")
    return nil, false
  }
  return msg, true
}
