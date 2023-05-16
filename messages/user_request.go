package messages

import (
  "time"
  "bytes"
  "crypto/sha512"
  "encoding/gob"
  "encoding/base64"
)

// Information about Organization
type Request struct {
  Id            string          `json:"id"`
  UserId        string          `json:"user_id"`
  Body          []byte          `json:"body"`
  CreatedAt     time.Time       `json:"created_at"`
  Sign          []byte          `json:"sign"`
}

func NewRequest() *Request {
  return &Request{}
}

func (i *Request) Init(id string, msg []byte) {
  i.Id = id
  i.Body = msg
  i.CreatedAt = time.Now()
}

func (i *Request) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write(i.Body)
  sha_512.Write([]byte(i.Id + i.CreatedAt.String()))
  return sha_512.Sum(nil)
}

func (i *Request) Pack() string {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return base64.StdEncoding.EncodeToString(buff.Bytes())
}

func (i *Request) Unpack(msg string) bool {
  input, err := base64.StdEncoding.DecodeString(msg)
  if err != nil {
    return false
  }
  buf := bytes.NewBuffer(input)
  decoder := gob.NewDecoder(buf)
  err = decoder.Decode(i)
  if err != nil {
    return false
  }
  return true
}
