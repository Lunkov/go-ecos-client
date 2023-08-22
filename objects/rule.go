package objects

import (
  "time"
  "github.com/google/uuid"
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  
  "go-ecos-client/utils"
)

type Action struct {
  Id           string    `db:"id"             json:"id,omitempty"              yaml:"id"`
  Func         string    `db:"func"           json:"func,omitempty"            yaml:"func"`

  Handler      string    `db:"handler"        json:"handler,omitempty"         yaml:"handler"`
  Description  string    `db:"description"    json:"description,omitempty"     yaml:"description"`
}

type LinkAction struct {
  Id           string    `db:"id"             json:"id,omitempty"              yaml:"id"`
  LinkFrom     string    `db:"link_from"      json:"link_from,omitempty"       yaml:"link_from"`
  LinkTo       string    `db:"link_to"        json:"link_to,omitempty"         yaml:"link_to"`
  Condition    string    `db:"condition"      json:"condition,omitempty"       yaml:"condition"`
  Description  string    `db:"description"    json:"description,omitempty"     yaml:"description"`
}

type Rule struct {
  Id           string
  TypeId       string
  ParentId     string
  CID          string
  Version      uint64
  Name         string
  Actions      map[string]Action
  LinkActions  map[string]LinkAction
  CreatedAt    time.Time
  UpdatedAt    time.Time
  DeletedAt   *time.Time
}

func (i *Rule) NewID() {
  uid, _ := uuid.NewUUID()
  i.Id = uid.String()
}

func (i *Rule) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.Id + i.TypeId + i.CreatedAt.String() + i.UpdatedAt.String()))
  sha_512.Write(utils.UInt64ToBytes(i.Version))
  return sha_512.Sum(nil)
}

func (i *Rule) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *Rule) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(i)
  if err != nil {
    return false
  }
  return true
}
