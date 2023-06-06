package client

import (
  "time"
  "math/rand"
  "net/url"
)

type ClientECOS struct {
  currentUrl       string
  currentIndexUrl  int
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
  client := &ClientECOS{url: au, maxRetries: maxRetries}
  client.selectRandomServer()
  return client
}

func (c *ClientECOS) selectRandomServer() {
  c.currentIndexUrl = rand.Intn(len(c.url))
  c.currentUrl = c.url[c.currentIndexUrl]
}

func (c *ClientECOS) getAnotherServer() {
  c.currentIndexUrl ++
  if len(c.url) <= c.currentIndexUrl {
    c.currentIndexUrl = 0
  }
  c.currentUrl = c.url[c.currentIndexUrl]
}
