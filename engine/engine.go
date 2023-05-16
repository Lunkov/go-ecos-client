package main

import (
  "encoding/json"
  "github.com/valyala/fasthttp"
  "github.com/valyala/fasthttprouter"
  "github.com/golang/glog"
  
  "github.com/Lunkov/go-ecos-client/messages"
)

type Engine struct {
}

func NewEngine() *Engine {
  a := &Engine{}
  return a
}

func (a *Engine) HandleObject(input []byte) ([]byte, bool) {
  inobj := messages.NewReqActionObject()
  if !inobj.Deserialize() {
    return nil, false
  }
  
}

func (a *Engine) HandleCoin(input []byte) ([]byte, bool) {
  inc := messages.NewReqActionCoin()
  if !inc.Deserialize() {
    return nil, false
  }
  inc.
  f, ok := [inc.IdAction]
}

func (a *Engine) HandleOrganization(input []byte) ([]byte, bool) {
  inobj := messages.NewReqActionObject()
  if !inobj.Deserialize() {
    return nil, false
  }
  
}

func (a *Engine) HandleUser(input []byte) ([]byte, bool) {
  inobj := messages.NewReqActionObject()
  if !inobj.Deserialize() {
    return nil, false
  }
  
}
