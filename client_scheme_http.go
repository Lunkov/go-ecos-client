package client

import (
  "time"
  "syscall"
  "errors"
  "net/http"
  "strings"
  "io/ioutil"
  
  "github.com/golang/glog"
)

func (c *ClientECOS) httpRequest(url string, request string) ([]byte, bool) {
  var answer []byte
  var reconnect, ok bool  
  retry := 1
  for {
    answer, reconnect, ok = c.httpReq("POST", url, request)
    if ok {
      return answer, true
    }
    if retry >= c.maxRetries {
      return nil, false
    }
    retry ++
    if reconnect {
      continue
    }
    if !ok {
      return nil, false
    }
  }
  return answer, true
}

func (c *ClientECOS) httpReq(protocol string, url string, request string) ([]byte, bool, bool) {
  client := http.Client{Timeout: time.Duration(2) * time.Second}
  req, errreq := http.NewRequest(protocol, c.currentUrl + url, strings.NewReader(request))
  if errreq != nil {
    glog.Errorf("ERR: request(%s): %v", url, errreq)
    return nil, true, false
  }
  req.Header.Add("Accept", `application/binary`)
  
  resp, errresp := client.Do(req)
  if errresp != nil {
    glog.Errorf("ERR: request.do(%s): %v", url, errresp, errors.Is(errresp, syscall.ECONNREFUSED))
    return nil, errors.Is(errresp, syscall.ECONNREFUSED), false
  }
  defer resp.Body.Close()

  answer, errbody := ioutil.ReadAll(resp.Body)
  if errbody != nil {
    glog.Errorf("ERR: request.read(%s): %v", url, errbody)
    return nil, false, false
  }
  
  return answer, false, true
}

