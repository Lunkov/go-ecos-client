package messages

import (
)

type TokenBlock struct {
  Name        string  `json:"name"`
  Symbol      string  `json:"symbol"`
  TotalSupply uint64  `json:"total_supply"`
  Decimals    uint8   `json:"decimals"`
  Owner       string  `json:"owner"`
}
