package objects

import (
  "time"
  "github.com/google/uuid"
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

type ContractSide struct {
  Role         string
  Address      string

  PublicKey   []byte          `json:"pubkey"`  
  Sign        []byte          `json:"sign"`
}

type ContractUserInput struct {
  Id           int32
  Title        string
  Role         string
  Inputs       map[string]string
}

type ContractUserInputData struct {
  Id           uint32
  Role         string
  CreatedAt    time.Time
  Inputs       map[string]string
}

type ContractSequence struct {
  Name         string

}

type ContractPayment struct {
  RoleFrom     string
  RoleTo       string
  
}

type Contract struct {
  Id           string
  Description  string
  
  Version      uint64
  TypeId       uint32

  UserInputs   []ContractUserInput

  Roles        []string
  
  Signs        map[string]ContractSide

  PrevHash     []byte
  PrevCID      string
  Action       string
  
  CreatedAt    time.Time
  UpdatedAt    time.Time
  DeletedAt   *time.Time
}

type ContractState struct {
  Id           string
  Description  string
  
  Version      uint64
  TypeId       uint32

  Signs        map[string]ContractSide

  PrevHash     []byte
  PrevCID      string
  Action       string
  
  CreatedAt    time.Time
  UpdatedAt    time.Time
  DeletedAt   *time.Time
}

func NewContract() *Contract {
  return &Contract{}
}

func (c *Contract) NewID() {
  uid, _ := uuid.NewUUID()
  c.Id = uid.String()
}

func (c *Contract) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(c.Id + c.Description + c.PrevCID + c.Action + c.CreatedAt.String() + c.UpdatedAt.String()))
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
