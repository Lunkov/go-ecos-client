package objects

import (
  "time"
  "github.com/google/uuid"
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

type Object struct {
  Id           string
  TypeId       string
  ParentId     string
  CID          string
  Version      uint64
  Data         []byte
  CreatedAt    time.Time
  UpdatedAt    time.Time
  DeletedAt   *time.Time
}

func (i *Object) NewID() {
  uid, _ := uuid.NewUUID()
  i.Id = uid.String()
}

func (i *Object) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.Id + i.TypeId + i.CreatedAt.String() + i.UpdatedAt.String()))
  sha_512.Write(utils.UInt64ToBytes(i.Version))
  sha_512.Write(i.Data)
  return sha_512.Sum(nil)
}

func (i *Object) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *Object) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(i)
  if err != nil {
    return false
  }
  return true
}
