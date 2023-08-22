package objects

import (
  "time"
  "github.com/google/uuid"
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  
  "go-ecos-client/utils"
)

const (
  ContractNew    = 0
  ContractActive = 100
  ContractClosed = 1000


  ContractPublic  = 2000
  ContractPrivate = 4000
  
  ContractTypePublicOffer  = 1
  ContractTypePrivateOffer = 2
  ContractTypePrivate      = 3
)

type ContractSide struct {
  Role         string
  Address      string

  PublicKey   []byte          `json:"pubkey"`  
  Sign        []byte          `json:"sign"`
}

type ContractRole struct {
  RoleId            string

  Title             string
  TitleTr           map[string]string

  Description       string
  DescriptionTr     map[string]string
}

/*
type ContractSequence struct {
  Name         string

}

type ContractPayment struct {
  RoleFrom     string
  RoleTo       string
  
}*/

type Contract struct {
  Id                string

  Description       string
  DescriptionTr     map[string]string

  Version           uint64
  TypeId            uint32
  
  DurationHours     uint32

  Actions           []ContractAction
  Payments          []ContractPayment

  UserForms         []UserForm

  Roles             []ContractRole
  
  Signs             map[string]ContractSide

  PrevHash          []byte
  PrevCID           string
  
  CreatedAt         time.Time
}

func NewContract() *Contract {
  return &Contract{
                   Actions: make([]ContractAction, 0),
                   UserForms: make([]UserForm, 0),
                   Payments: make([]ContractPayment, 0),
                   Roles: make([]ContractRole, 0),
                   Signs: make(map[string]ContractSide),
                   }
}

func (c *Contract) NewID() {
  uid, _ := uuid.NewUUID()
  c.Id = uid.String()
}

func (c *Contract) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(c.Id + c.Description + c.PrevCID + c.CreatedAt.String()))
  sha_512.Write(utils.UInt64ToBytes(c.Version))
  sha_512.Write(utils.UInt32ToBytes(c.TypeId))
  sha_512.Write(c.PrevHash)
  return sha_512.Sum(nil)
}

func (c *Contract) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(c)
  return buff.Bytes()
}

func (c *Contract) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(c)
  if err != nil {
    return false
  }
  return true
}

func (c *Contract) GetRolesList() []string {
  result := make([]string, 0)
  for _, role := range c.Roles {
    result = append(result, role.RoleId)
  }
  return result
}

func (c *Contract) GetActionsList() []string {
  result := make([]string, 0)
  for _, action := range c.Actions {
    result = append(result, action.ActionId)
  }
  return result
}

func (c *Contract) GetFormsList() []string {
  result := make([]string, 0)
  for _, form := range c.UserForms {
    result = append(result, form.FormId)
  }
  return result
}
