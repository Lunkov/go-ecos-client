package client

import (
  "errors"
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"

  "go-ecos-client/objects"
  "go-ecos-client/messages"
)

func (c *ClientECOS) GetBalance(w wallets.IWallet) (*messages.Balance, error) {
  msg := messages.NewGetBalanceReq()
  errs := msg.Init(w, hdwallet.ECOS)
  if errs != nil {
    return nil, errs
  }
  
  answer, err := c.httpRequest("POST", "/wallet/balance", string(msg.Serialize()))
  if err != nil {
    return nil, err
  }
  result := messages.NewBalance()
  err = result.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return result, nil
}

func (c *ClientECOS) GetWalletTX(w wallets.IWallet) (*objects.WalletTransactions, error) {
  msg := messages.NewGetBalanceReq()
  errs := msg.Init(w, hdwallet.ECOS)
  if errs != nil {
    return nil, errs
  }
  
  answer, err := c.httpRequest("POST", "/wallet/tx", string(msg.Serialize()))
  if err != nil {
    return nil, err
  }
  result := objects.NewWalletTransactionsEmpty()
  err = result.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return result, nil
}

func (c *ClientECOS) TransactionNew(w wallets.IWallet, addressTo string, coin uint32, value uint64) (*objects.Transaction, error) {
  msg := messages.NewMsgTransaction()
  msg.Init(messages.StatusTxNew, w, addressTo, coin, coin, 0, value)
  
  errs := msg.DoSign(w)
  if errs != nil {
    return nil, errs
  }
  
  answer, err := c.httpRequest("POST", "/transaction/new", string(msg.Serialize()))
  if err != nil {
    return nil, err
  }

  msgAnswer := objects.NewTX()
  err = msgAnswer.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return msgAnswer, nil
}

func (c *ClientECOS) TransactionStatus(w wallets.IWallet, IdTx []byte) (*messages.MsgTransactionStatus, error) {
  msg := messages.NewMsgTransactionStatus(IdTx)
  
  err := msg.DoSign(w)
  if err != nil {
    return nil, err
  }
  
  answer, errc := c.httpRequest("POST", "/transaction/status", string(msg.Serialize()))
  if errc != nil {
    return nil, errc
  }

  err = msg.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return msg, nil
}

func (c *ClientECOS) TransactionCommit(w wallets.IWallet, tx *objects.Transaction) (*objects.Transaction, error) {
  if tx == nil {
    return nil, errors.New("TX is empty")
  }
  
  msg := tx
  
  err := msg.DoSign(w)
  if err != nil {
    return nil, err
  }
  output, errs := msg.Serialize()
  if errs != nil {
    return nil, errs
  }
  
  answer, errc := c.httpRequest("POST", "/transaction/commit", string(output))
  if errc != nil {
    return nil, errc
  }
  err = msg.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return msg, nil
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
    return nil, false
  }
  return result, true
}
*/
