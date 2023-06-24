package objects

import (
)

const (
  FormInputTypeUndef     = 0
  FormInputTypeString    = 10
  FormInputTypePassword  = 11
  FormInputTypePassword2 = 12
  FormInputTypeText      = 15
  FormInputTypeNumber    = 20
  FormInputTypeCheckBox  = 30
  FormInputTypeFileImage = 100
  FormInputTypeFileDoc   = 101
)

type UserFormInput struct {
  InputId      string
  GroupId      string
  TypeId       uint32
  Required     bool
  Title        string
  TitleTr      map[string]string
}

type UserFormGroup struct {
  GroupId      string
  Title        string
  TitleTr      map[string]string
}

type UserForm struct {
  FormId            string

  ParentFormId      string

  Title             string
  TitleTr           map[string]string

  Description       string
  DescriptionTr     map[string]string

  Groups       []UserFormGroup
  Inputs       []UserFormInput
}
