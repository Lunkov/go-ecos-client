package messages

import (
)

// Information about Organization
type OrgInfo struct {
  Name          string          `json:"name"`
  Logo          string          `json:"logo"`
  URL           string          `json:"url"`
  MainCoin      string          `json:"maincoin"`
}

func NewOrgInfo() *OrgInfo {
  return &OrgInfo{}
}
