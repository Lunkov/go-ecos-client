package messages

import (
  "time"
  "bytes"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/gob"
  "encoding/binary"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

type ReqActionObject struct {
  Version       string          `json:"version"`
  
  IdMessage     uint32          `json:"id_message"`
  
  IdObject      string          `json:"id_object"`
  IdAction      uint32          `json:"id_action"`
  DataObject    []byte          `json:"data_object"`

  UpdatedAt     time.Time       `json:"updated_at"`

  PubKey        []byte          `json:"pubkey"`
  Sign          []byte          `json:"sign"`
}

func NewReqActionObject() *ReqActionObject {
  return &ReqActionObject{Version: "1"}
}

func (i *ReqActionObject) Init(idAction uint32, idObject string, dataObject []byte) {
  i.IdAction = idAction
  i.IdObject = idObject
  i.DataObject = dataObject
  i.UpdatedAt = time.Now()
  i.IdMessage = i.msgId()
}

func (i *ReqActionObject) msgId() uint32 {
  sha_256 := sha256.New()
  sha_256.Write([]byte(i.IdObject))
  sha_256.Write(utils.UInt32ToBytes(i.IdAction))
  sha_256.Write(i.DataObject)
  sha_256.Write([]byte(i.UpdatedAt.String()))
  return binary.LittleEndian.Uint32(sha_256.Sum(nil))
}

func (i *ReqActionObject) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.Version + i.IdObject))
  sha_512.Write(utils.UInt32ToBytes(i.IdAction))
  sha_512.Write(utils.UInt32ToBytes(i.IdMessage))
  sha_512.Write(i.DataObject)
  return sha_512.Sum(nil)
}


func (i *ReqActionObject) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *ReqActionObject) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(i) == nil
}
