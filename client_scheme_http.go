package client

import (
  "time"
  "syscall"
  "errors"
  "net/http"
  "strings"
  "io"
)

func (c *ClientECOS) httpRequest(protocol string, url string, request string) ([]byte, error) {
  var answer []byte
  var reconnect bool
  var err error 
  retry := 1
  for {
    answer, reconnect, err = c.httpReq(protocol, url, request)
    if err == nil {
      return answer, nil
    }
    c.getAnotherServer()
    if retry >= c.maxRetries {
      return nil, errors.New("HTTP maxRetries")
    }
    retry ++
    if reconnect {
      continue
    }
    if err != nil {
      return nil, err
    }
  }
  return answer, nil
}

func (c *ClientECOS) httpReq(protocol string, url string, request string) ([]byte, bool, error) {
  client := http.Client{Timeout: time.Duration(2) * time.Second}
  req, errreq := http.NewRequest(protocol, c.currentUrl + url, strings.NewReader(request))
  if errreq != nil {
    return nil, true, errreq
  }
  req.Header.Add("Accept", `application/binary`)
  
  resp, errresp := client.Do(req)
  if errresp != nil {
    return nil, errors.Is(errresp, syscall.ECONNREFUSED), errresp
  }
  defer resp.Body.Close()
  
  answer, errbody := io.ReadAll(resp.Body)
  if errbody != nil {
    return nil, false, errbody
  }
  
  return answer, false, nil
}

