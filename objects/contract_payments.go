package objects

import (
)

type ContractPayment struct {
  FinishActionId  string
  
  AddresFrom      string
  AddresTo        string

  Costs           uint64
}
