package objects

import (
)

type ContractPayment struct {
  Title               string
  TitleTr             map[string]string

  Description         string
  DescriptionTr       map[string]string

  FinishActionId      string
  
  RoleAddresFrom      string
  RoleAddresTo        string

  Amount              uint64
}
