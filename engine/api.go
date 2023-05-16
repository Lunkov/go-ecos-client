package main

import (
  "encoding/json"
  "github.com/valyala/fasthttp"
  "github.com/valyala/fasthttprouter"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-messages"
)

type API struct {
}

func NewAPI() *API {
  a := &API{}
  return a
}

func (a *API) GetBalance(inputMessage []byte) ([]byte, bool) {
}

func (a *API) GetBalance(inputMessage []byte) ([]byte, bool) {
  ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")
  ctx.Response.Header.Set("Content-Type", "application/x-binary")
  
  req := messages.NewReqGetBalance()
  
  if !req.Deserialize(ctx.PostBody()) {
    ctx.SetStatusCode(fasthttp.StatusBadRequest)
    return
  }
  if glog.V(9) {
    glog.Infof("DBG: GetBalance(%s)", req.Address)
  }
  balance, ok := a.BPool.GetBalance(req.Address)
  if !ok {
    ctx.SetStatusCode(fasthttp.StatusNotFound)
    return 
  }
  ctx.SetStatusCode(fasthttp.StatusOK)
  ctx.Write(balance.Serialize())
}

func (a *API) NewTransaction(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
  ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST")
  ctx.Response.Header.Set("Content-Type", "application/x-binary")
  
  req := messages.NewTokenTransaction()
  
  if !req.Deserialize(ctx.PostBody()) {
    ctx.SetStatusCode(fasthttp.StatusBadRequest)
    return
  }
  if glog.V(9) {
    glog.Infof("DBG: REGISTER: Date (%v)", req)
  }
  
  obj := a.BPool.AddTransaction(req)
  jsonAnswer, err := json.Marshal(obj)
  if err != nil {
    ctx.SetStatusCode(fasthttp.StatusInternalServerError)
    return 
  }
  ctx.Write(jsonAnswer)
  ctx.SetStatusCode(fasthttp.StatusOK)
}

func (a *API) GetLastBlock(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
  ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")
  ctx.Response.Header.Set("Content-Type", "application/json")

  if glog.V(9) {
    glog.Infof("DBG: GetLastBlock")
  }
  obj := a.BPool.GetLastBlock()
  jsonAnswer, err := json.Marshal(obj)
  if err != nil {
    ctx.SetStatusCode(fasthttp.StatusInternalServerError)
    return 
  }

  ctx.SetStatusCode(fasthttp.StatusOK)
  ctx.Write(jsonAnswer)
}

func (a *API) GetLastCID(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
  ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")
  ctx.Response.Header.Set("Content-Type", "application/text")
  
  if glog.V(9) {
    glog.Infof("DBG: GetLastCID")
  }

  ctx.Write([]byte(a.BPool.GetLastCID()))
  ctx.SetStatusCode(fasthttp.StatusOK)
}

func (a *API) GetStat(ctx *fasthttp.RequestCtx, params fasthttprouter.Params) {
  ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET")
  ctx.Response.Header.Set("Content-Type", "application/json")

  if glog.V(9) {
    glog.Infof("DBG: GetLastBlock")
  }
  obj := a.BPool.GetStat()
  jsonAnswer, err := json.Marshal(obj)
  if err != nil {
    ctx.SetStatusCode(fasthttp.StatusInternalServerError)
    return 
  }

  ctx.SetStatusCode(fasthttp.StatusOK)
  ctx.Write(jsonAnswer)
}
