package objects

import (
)

type ContractActionStep struct {
  StepId          string

  Description     string
  DescriptionTr   map[string]string

  
  
}

type ContractAction struct {
  ActionId        string
  
  Description     string
  DescriptionTr   map[string]string

  UserAction      bool             
  RoleId          string

  Timeout         bool

  Steps          []ContractActionStep
}
