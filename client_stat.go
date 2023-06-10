package client

import (
  "github.com/Lunkov/go-ecos-client/messages"
)

func (c *ClientECOS) GetNetworkStat() (*messages.NetworkStat, bool) {
  answer, ok := c.httpRequest("POST", "/network/stat", "")
  if !ok {
    return nil, false
  }
  result := messages.NewNetworkStat()
  if !result.Deserialize(answer) {
    return nil, false
  }
  return result, true
}

func (c *ClientECOS) GetNodeStat() (*messages.NodeStat, bool) {
  answer, ok := c.httpRequest("POST", "/node/stat", "")
  if !ok {
    return nil, false
  }

  msgAnswer := messages.NewNodeStat()
  if !msgAnswer.Deserialize(answer) {
    return nil, false
  }
  return msgAnswer, true
}
