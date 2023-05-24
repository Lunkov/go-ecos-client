package client

import (
  "time"
  "bytes"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/gob"
  "encoding/binary"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

type Msg struct {
  Version       string          `json:"version"`
  
  IdMessage     uint32          `json:"id_message"`
  IdObject      string          `json:"id_object"`
  IdAction      uint32          `json:"id_action"`

  UpdatedAt     time.Time       `json:"updated_at"`
  PubKey        []byte          `json:"pubkey"`
  Sign          []byte          `json:"sign"`

  DataHash func() []byte        `json:"-"`
}

func NewMsg() *Msg {
  return &Msg{Version: "1"}
}

func (i *Msg) Init(idAction uint32, idObject string) {
  i.IdAction = idAction
  i.IdObject = idObject
  i.UpdatedAt = time.Now()
  i.IdMessage = i.msgId()
}

func (i *Msg) msgId() uint32 {
  sha_256 := sha256.New()
  sha_256.Write([]byte(i.IdObject))
  sha_256.Write(utils.UInt32ToBytes(i.IdAction))
  if i.DataHash != nil  {
    sha_256.Write(i.DataHash())
  }
  sha_256.Write([]byte(i.UpdatedAt.String()))
  return binary.LittleEndian.Uint32(sha_256.Sum(nil))
}

func (i *Msg) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.Version + i.IdObject))
  sha_512.Write(utils.UInt32ToBytes(i.IdAction))
  sha_512.Write(utils.UInt32ToBytes(i.IdMessage))
  if i.DataHash != nil {
    sha_512.Write(i.DataHash())
  }
  return sha_512.Sum(nil)
}


func (i *Msg) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *Msg) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(i) == nil
}
