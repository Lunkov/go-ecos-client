package client

import (
  "github.com/Lunkov/go-ecos-client/messages"
)

func (c *ClientECOS) GetNetworkStat() (*messages.NetworkStat, error) {
  answer, err := c.httpRequest("POST", "/network/stat", "")
  if err != nil {
    return nil, err
  }
  result := messages.NewNetworkStat()
  err = result.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return result, nil
}

func (c *ClientECOS) GetNodeStat() (*messages.NodeStat, error) {
  answer, err := c.httpRequest("POST", "/node/stat", "")
  if err != nil {
    return nil, err
  }

  msgAnswer := messages.NewNodeStat()
  err = msgAnswer.Deserialize(answer)
  if err != nil {
    return nil, err
  }
  return msgAnswer, nil
}
