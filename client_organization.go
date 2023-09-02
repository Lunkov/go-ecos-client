package client

import (
  "errors"
  "github.com/Lunkov/lib-wallets"

  "github.com/Lunkov/go-ecos-client/objects"
  "github.com/Lunkov/go-ecos-client/messages"
)

func (c *ClientECOS) OrganizationStatus(w wallets.IWallet, IdTx []byte) (*messages.MsgTransactionStatus, error) {
  msg := messages.NewMsgTransactionStatus(IdTx)
  
  errs := msg.DoSign(w)
  if errs != nil {
    return nil, errs
  }
  
  answer, err := c.httpRequest("POST", "/organization/status", string(msg.Serialize()))
  if err != nil {
    return nil, err
  }
  err = msg.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return msg, nil
}

func (c *ClientECOS) OrganizationNew(w wallets.IWallet, coin uint32) (*objects.Transaction, error) {
  msg := messages.NewMsgTransaction()
  msg.InitNFT(messages.StatusTxNewNFT, w, coin)
  
  errs := msg.DoSign(w)
  if errs != nil {
    return nil, errs
  }
  
  answer, err := c.httpRequest("POST", "/organization/new", string(msg.Serialize()))
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

func (c *ClientECOS) OrganizationCommit(w wallets.IWallet, coin uint32, org *messages.Organization) (*objects.Transaction, error) {
  if org == nil {
    return nil, errors.New("Organization is empty")
  }
  errs := org.DoSign(w)
  if errs != nil {
    return nil, errs
  }
  output, errs2 := org.Serialize()
  if errs2 != nil {
    return nil, errs2
  }
  
  answerBuf, err := c.httpRequest("POST", "/organization/commit", string(output))
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

func (c *ClientECOS) OrganizationGet(cid string) (*objects.OrganizationInfo, error) {
  
  answerBuf, err := c.httpRequest("POST", "/object", string(cid))
  if err != nil {
    return nil, err
  }
  // DECRYPT
  org := objects.NewOrganizationInfo()
  
  err = org.Deserialize(answerBuf)
  if err != nil {
    return nil, err
  }
  return org, nil
}
