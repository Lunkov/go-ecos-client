package client

import (
  "time"
  "math/rand"
  "net/url"
)

type Client struct {
  currentUrl       string
  currentProtocol  string
  url            []string
  maxRetries       int
}

func NewClient(urls []string, maxRetries int) (*Client) {
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
  return &Client{url: au, maxRetries: maxRetries}
}

func (c *Client) selectServer() {
  index := rand.Intn(len(c.url))
  c.currentUrl = c.url[index]
}

