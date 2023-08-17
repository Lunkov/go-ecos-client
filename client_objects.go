package client

import (
  /*
  "github.com/Lunkov/lib-wallets"
  
  "github.com/Lunkov/go-ecos-client/objects"
  */
)
/*
func (c *Client) GetObject(id string) (*objects.Object, bool) {
  c.selectServer()
  msg := messages.ReqActionObject()
  msg.Address = w.GetAddress("ECOS")
  
  answer, ok := c.httpRequest("/object/", string(msg.Serialize()))
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

func (c *Client) GetObjectList(page int, pageSize int) (*objects.Object, int, bool) {
  c.selectServer()
  msg := messages.NewReqGetBalance()
  msg.Address = w.GetAddress("ECOS")
  
  answer, ok := c.httpRequest("/objects", string(msg.Serialize()))
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

func (c *Client) PutObject(obj *objects.Object) (*objects.Object, bool) {
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
}*/
