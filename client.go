package client

import (
  "time"
  "math/rand"
  "net/url"
)

type ClientECOS struct {
  currentUrl       string
  currentProtocol  string
  url            []string
  maxRetries       int
}

func NewClientECOS(urls []string, maxRetries int) (*ClientECOS) {
  rand.Seed(time.Now().UnixNano())
  au := make([]string, 0)
  for _, ui := range urls {
    u, err := url.Parse(ui)
    if err != nil {
      continue
    }
    if u.Scheme == "http" || u.Scheme == "https" {
      au = append(au, ui)
    }
  }
  return &ClientECOS{url: au, maxRetries: maxRetries}
}

func (c *ClientECOS) selectServer() {
  index := rand.Intn(len(c.url))
  c.currentUrl = c.url[index]
}

