package client

import (
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"

  "github.com/Lunkov/go-ecos-client/messages"
)

func (c *ClientECOS) GetBalance(w wallets.IWallet) (*messages.Balance, bool) {
  msg := messages.NewGetBalanceReq()
  oks := msg.Init(w, hdwallet.ECOS)
  if !oks {
    return nil, false
  }
  
  answer, ok := c.httpRequest("POST", "/wallet/balance", string(msg.Serialize()))
  if !ok {
    return nil, false
  }
  result := messages.NewBalance()
  if !result.Deserialize(answer) {
    return nil, false
  }
  return result, true
}

func (c *ClientECOS) TransactionNew(w wallets.IWallet, addressTo string, coin uint32, value uint64) (*messages.Transaction, bool) {
  msg := messages.NewMsgTransaction()
  msg.Init(messages.StatusTxNew, w, addressTo, coin, coin, 0, value)
  
  if !msg.DoSign(w) {
    return nil, false
  }
  
  answer, ok := c.httpRequest("POST", "/transaction/new", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  msgAnswer := messages.NewTX()
  if !msgAnswer.Deserialize(answer) {
    return nil, false
  }
  return msgAnswer, true
}

func (c *ClientECOS) TransactionStatus(w wallets.IWallet, IdTx []byte) (*messages.MsgTransactionStatus, bool) {
  msg := messages.NewMsgTransactionStatus(IdTx)
  
  if !msg.DoSign(w) {
    return nil, false
  }
  
  answer, ok := c.httpRequest("POST", "/transaction/status", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  if !msg.Deserialize(answer) {
    return nil, false
  }
  return msg, true
}

func (c *ClientECOS) TransactionCommit(w wallets.IWallet, tx *messages.Transaction) (*messages.Transaction, bool) {
  if tx == nil {
    return nil, false
  }
  
  msg := tx
  
  if !msg.DoSign(w) {
    return nil, false
  }
  output, oko := msg.Serialize()
  if !oko {
    return nil, false
  }
  
  answer, ok := c.httpRequest("POST", "/transaction/commit", string(output))
  if !ok {
    return nil, false
  }

  if !msg.Deserialize(answer) {
    return nil, false
  }
  return msg, true
}

/*
func (c *ClientECOS) GetNodeStat() (*messages.Balance, bool) {
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
    if glog.V(2) {
      glog.Errorf("ERR: GetBalance.Deserialize")
    }
    return nil, false
  }
  return result, true
}

func (c *ClientECOS) GetNetworkStat() (*messages.Balance, bool) {
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
    if glog.V(2) {
      glog.Errorf("ERR: GetBalance.Deserialize")
    }
    return nil, false
  }
  return result, true
}
*/
