package client

import (
  "github.com/Lunkov/go-ecos-client/objects"
  //"github.com/Lunkov/go-ecos-client/messages"
)

func (c *ClientECOS) SignGet(cid string) (*objects.SignInfo, error) {
  
  answerBuf, err := c.httpRequest("POST", "/object", string(cid))
  if err != nil {
    return nil, err
  }
  // DECRYPT
  sign := objects.NewSignInfo()
  
  err = sign.Deserialize(answerBuf)
  if err != nil {
    return nil, err
  }
  return sign, nil
}
