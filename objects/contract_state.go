package objects

import (
  "time"
  "github.com/google/uuid"
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  
  "go-ecos-client/utils"
)

type ContractState struct {
  Id           string
  Description  string

  ActionId     string
  StepId       string
  
  Version      uint64
  TypeId       uint32
  StatusId     uint32

  Data         []byte
  
  Signs        map[string]ContractSide

  PrevHash     []byte
  PrevCID      string
  
  CreatedAt    time.Time
}

func NewContractState() *ContractState {
  return &ContractState{}
}

func (c *ContractState) NewID() {
  uid, _ := uuid.NewUUID()
  c.Id = uid.String()
}

func (c *ContractState) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(c.Id + c.Description + c.CreatedAt.String()))
  sha_512.Write(utils.UInt64ToBytes(c.Version))
  sha_512.Write(utils.UInt32ToBytes(c.TypeId))
  sha_512.Write(c.PrevHash)
  return sha_512.Sum(nil)
}

func (c *ContractState) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(c)
  return buff.Bytes()
}

func (c *ContractState) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(c)
  if err != nil {
    return false
  }
  return true
}
