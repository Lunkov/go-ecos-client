package objects

import (
)

const (
  ActionUser        = 1
  ActionTimer       = 2
  ActionCircleTimer = 3
  ActionEvent       = 4
)

type ContractActionLinkStep struct {
  FromStepId      string
  toStepId        string

  Description     string
  DescriptionTr   map[string]string

  Condition       string
}

type ContractActionStep struct {
  StepId          string

  Description     string
  DescriptionTr   map[string]string

  UserFormId      string
  
  ServiceId       string

  SendEventId     string
  
  NextSteps       []ContractActionLinkStep
}

type ContractAction struct {
  ActionId        string
  TypeActionId    uint32
  
  Description     string
  DescriptionTr   map[string]string

  RoleId          string

  ReceiveEventId  string

  TypeTimer       uint64
  Timer           string

  Steps          []ContractActionStep
}
