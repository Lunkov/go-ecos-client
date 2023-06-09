package client

import (
  "github.com/Lunkov/lib-wallets"

  "github.com/Lunkov/go-ecos-client/objects"
  "github.com/Lunkov/go-ecos-client/messages"
)

func (c *ClientECOS) PassportStatus(w wallets.IWallet, IdTx []byte) (*messages.MsgTransactionStatus, bool) {
  msg := messages.NewMsgTransactionStatus(IdTx)
  
  if !msg.DoSign(w) {
    return nil, false
  }
  
  answer, ok := c.httpRequest("POST", "/passport/status", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  if !msg.Deserialize(answer) {
    return nil, false
  }
  return msg, true
}

func (c *ClientECOS) PassportNew(w wallets.IWallet, coin uint32) (*objects.Transaction, bool) {
  msg := messages.NewMsgTransaction()
  msg.InitNFT(messages.StatusTxNewNFT, w, coin)
  
  if !msg.DoSign(w) {
    return nil, false
  }
  
  answer, ok := c.httpRequest("POST", "/passport/new", string(msg.Serialize()))
  if !ok {
    return nil, false
  }

  msgAnswer := objects.NewTX()
  if !msgAnswer.Deserialize(answer) {
    return nil, false
  }
  return msgAnswer, true
}

func (c *ClientECOS) PassportCommit(w wallets.IWallet, coin uint32, passport *messages.PassportInfo) (*objects.Transaction, bool) {
  if passport == nil {
    return nil, false
  }
  if !passport.Passport.DoSignPerson(w) {
    return nil, false
  }
  output, oko := passport.Serialize()
  if !oko {
    return nil, false
  }
  
  answerBuf, ok := c.httpRequest("POST", "/passport/commit", string(output))
  if !ok {
    return nil, false
  }
  answer := objects.NewTX()
  if !answer.Deserialize(answerBuf) {
    return nil, false
  }
  return answer, true
}

func (c *ClientECOS) PassportGet(cid string, password string) (*objects.Passport, bool) {
  
  answerBuf, ok := c.httpRequest("POST", "/object", string(cid))
  if !ok {
    return nil, false
  }
  // DECRYPT
  passport := objects.NewPassport()
  
  if !passport.DeserializeDecrypt(password, answerBuf) {
    return nil, false
  }
  return passport, true
}
