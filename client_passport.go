package client

import (
  "errors"
  "github.com/Lunkov/lib-wallets"

  "go-ecos-client/objects"
  "go-ecos-client/messages"
)

func (c *ClientECOS) PassportStatus(w wallets.IWallet, IdTx []byte) (*messages.MsgTransactionStatus, error) {
  msg := messages.NewMsgTransactionStatus(IdTx)
  
  errs := msg.DoSign(w)
  if errs != nil {
    return nil, errs
  }
  
  answer, err := c.httpRequest("POST", "/passport/status", string(msg.Serialize()))
  if err != nil {
    return nil, err
  }
  err = msg.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return msg, nil
}

func (c *ClientECOS) PassportNew(w wallets.IWallet, coin uint32) (*objects.Transaction, error) {
  msg := messages.NewMsgTransaction()
  msg.InitNFT(messages.StatusTxNewNFT, w, coin)
  
  errs := msg.DoSign(w)
  if errs != nil {
    return nil, errs
  }
  
  answer, err := c.httpRequest("POST", "/passport/new", string(msg.Serialize()))
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

func (c *ClientECOS) PassportCommit(w wallets.IWallet, coin uint32, passport *messages.PassportInfo) (*objects.Transaction, error) {
  if passport == nil {
    return nil, errors.New("Passport is empty")
  }
  errs := passport.Passport.DoSignPerson(w)
  if errs != nil {
    return nil, errs
  }
  output, errs2 := passport.Serialize()
  if errs2 != nil {
    return nil, errs2
  }
  
  answerBuf, err := c.httpRequest("POST", "/passport/commit", string(output))
  if err != nil {
    return nil, err
  }
  answer := objects.NewTX()
  err = answer.Deserialize(answerBuf)
  if err != nil {
    return nil, err
  }
  return answer, nil
}

func (c *ClientECOS) PassportGet(cid string, password string) (*objects.Passport, error) {
  
  answerBuf, err := c.httpRequest("POST", "/object", string(cid))
  if err != nil {
    return nil, err
  }
  // DECRYPT
  passport := objects.NewPassport()
  
  err = passport.DeserializeDecrypt(password, answerBuf)
  if err != nil {
    return nil, err
  }
  return passport, nil
}
