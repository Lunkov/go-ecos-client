package messages

import (
  "time"
  "bytes"
  
  "crypto/ecdsa"

  
  "crypto/sha512"
  
  "encoding/gob"
  "encoding/binary"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

const (
  GetBalance = 0
  CalcTransaction = 1
  DoTransaction = 2
  
  
  StatusNew = 0
  StatusProcessing = 100
  StatusDone = 10000
)

type ReqActionCoin struct {
  Version       string          `json:"version"`
  
  IdMessage     uint64          `json:"id_message"`
  IdStatus      uint64          `json:"id_status"`
  
  IdAction      uint32          `json:"id_action"`
  
  AddressFrom   string          `json:"address_from"`
  AddressTo     string          `json:"address_to"`
  CoinFrom      string          `json:"coin_from"`
  CoinTo        string          `json:"coin_to"`
  ValueFrom     uint64          `json:"value"`
  ValueTo       uint64          `json:"value"`

  UpdatedAt     time.Time       `json:"updated_at"`

  PubKey        []byte          `json:"pubkey"`
  Sign          []byte          `json:"sign"`
}

func NewReqActionCoin() *ReqActionCoin {
  return &ReqActionCoin{Version: "1"}
}

func (i *ReqActionCoin) Init(idAction uint32,
                               addressFrom string,
                               addressTo string, 
                               coinFrom string, 
                               coinTo string, 
                               valueFrom uint64,
                               valueTo uint64,
                               ) {
  i.IdAction = idAction
  i.AddressFrom = addressFrom
  i.AddressTo = addressTo
  i.CoinFrom = coinFrom
  i.CoinTo = coinTo
  i.ValueFrom = valueFrom
  i.ValueTo = valueTo
  i.UpdatedAt = time.Now()
  i.IdMessage = i.msgId()
}

func (i *ReqActionCoin) msgId() uint64 {
  sha := sha512.New()
  sha.Write([]byte(i.AddressFrom))
  sha.Write([]byte(i.AddressTo))
  sha.Write([]byte(i.CoinFrom))
  sha.Write([]byte(i.CoinTo))
  sha.Write(utils.UInt64ToBytes(i.ValueFrom))
  sha.Write(utils.UInt64ToBytes(i.ValueTo))
  sha.Write(utils.UInt32ToBytes(i.IdAction))
  sha.Write([]byte(i.UpdatedAt.String()))
  return binary.LittleEndian.Uint64(sha.Sum(nil))
}

func (i *ReqActionCoin) Hash() []byte {
  sha := sha512.New()
  sha.Write([]byte(i.Version))
  sha.Write(utils.UInt32ToBytes(i.IdAction))
  sha.Write(utils.UInt64ToBytes(i.IdMessage))
  sha.Write(utils.UInt64ToBytes(i.IdStatus))
  sha.Write([]byte(i.AddressFrom))
  sha.Write([]byte(i.AddressTo))
  sha.Write([]byte(i.CoinFrom))
  sha.Write([]byte(i.CoinTo))
  sha.Write(utils.UInt64ToBytes(i.ValueFrom))
  sha.Write(utils.UInt64ToBytes(i.ValueTo))
  sha.Write([]byte(i.UpdatedAt.String()))
  return sha.Sum(nil)
}


func (i *ReqActionCoin) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *ReqActionCoin) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(i) == nil
}

func (i *ReqActionCoin) DoSign(pk *ecdsa.PrivateKey) bool {
  if i.AddressFrom != utils.PubkeyToAddress(&pk.PublicKey).Hex() {
    return false
  }
  sign, ok := utils.ECDSA256SignHash512(pk, i.Hash())
  if !ok {
    return false
  }
  i.Sign = sign
  i.PubKey, ok = utils.ECDSAPublicKeySerialize(&pk.PublicKey)
  if !ok {
    return false
  }
  return true
}

func (i *ReqActionCoin) DoVerify() bool {
  return utils.ECDSA256VerifySender(i.AddressFrom, i.PubKey, i.Hash(), i.Sign)
}
