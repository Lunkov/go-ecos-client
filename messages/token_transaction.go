package messages

import (
  "bytes"
  "strconv"
  
  "crypto/sha512"
  
  "encoding/gob"
)

type TokenTransaction struct {
  AddressFrom        string  `json:"address_from"`
  AddressTo          string  `json:"address_to"`
  Coin               uint32  `json:"coin"`
  Value              uint64  `json:"value"`
  MaxCost            uint64  `json:"max_cost"`
  PublicKey          []byte  `json:"public_key"`
  Sign               []byte  `json:"sign"`
}

func NewTokenTransaction() (*TokenTransaction) {
  return &TokenTransaction{}
}

func (t *TokenTransaction) Init(addressFrom string, addressTo string, coin uint32, value uint64, maxCost uint64, publicKey []byte) {
  t.AddressFrom = addressFrom
  t.AddressTo = addressTo
  t.Coin = coin
  t.Value = value
  t.MaxCost = maxCost
  t.PublicKey = publicKey
}

func (t *TokenTransaction) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(t.AddressFrom + t.AddressTo + strconv.FormatUint(uint64(t.Coin), 10) + strconv.FormatUint(t.Value, 10) + strconv.FormatUint(t.MaxCost, 10) + string(t.PublicKey)))
  return sha_512.Sum(nil)
}

func (t *TokenTransaction) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(t)
  return buff.Bytes()
}

func (t *TokenTransaction) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(t)
  if err != nil {
    return false
  }
  return true
}
