package client

import (
  "bytes"
  "time"
  "crypto/sha512"
  "encoding/gob"
  "encoding/json"
  
  "github.com/Lunkov/lib-wallets"
  "github.com/Lunkov/go-ecos-client/utils"
)

const (
  StatusTxNew       = 0
  StatusTxVerifying = 2
  StatusTxVerified  = 10
  StatusTxApproved  = 100
  StatusTxSaved     = 500
)

type MsgTransaction struct {
  Msg

  AddressFrom   string          `json:"address_from"`
  AddressTo     string          `json:"address_to"`
  CoinFrom      uint32          `json:"coin_from"`
  CoinTo        uint32          `json:"coin_to"`
  ValueFrom     uint64          `json:"value"`
  ValueTo       uint64          `json:"value"`
  
}

func NewMsgTransaction() *MsgTransaction {
  t := &MsgTransaction{}
  t.DataHash = t.HashMsg
  return t
}


func (i *MsgTransaction) Init(idAction uint32,
                               wallet wallets.IWallet,
                               addressTo string, 
                               coinFrom uint32, 
                               coinTo uint32, 
                               valueFrom uint64,
                               valueTo uint64,
                               ) {
  i.IdAction = idAction
  i.IdObject = "payment"
  i.AddressFrom = wallet.GetAddress(coinFrom)
  i.AddressTo = addressTo
  i.CoinFrom = coinFrom
  i.CoinTo = coinTo
  i.ValueFrom = valueFrom
  i.ValueTo = valueTo
  i.UpdatedAt = time.Now()
  i.IdMessage = i.msgId()
}

func (i *MsgTransaction) HashMsg() []byte {
  sha := sha512.New()
  sha.Write([]byte(i.AddressFrom))
  sha.Write([]byte(i.AddressTo))
  sha.Write(utils.UInt32ToBytes(i.CoinFrom))
  sha.Write(utils.UInt32ToBytes(i.CoinTo))
  sha.Write(utils.UInt64ToBytes(i.ValueFrom))
  sha.Write(utils.UInt64ToBytes(i.ValueTo))
  sha.Write(utils.UInt32ToBytes(i.IdAction))
  
  return sha.Sum(nil)
}


func (i *MsgTransaction) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *MsgTransaction) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(i) == nil
}

func (i *MsgTransaction) ToJSON() (string, bool) {
  jsonAnswer, err := json.Marshal(i)
  if err != nil {
    return "", false
  }
  return string(jsonAnswer), true
}

func (i *MsgTransaction) FromJSON(msg string) bool {
  if err := json.Unmarshal([]byte(msg), i); err != nil {
    return false
  }
  return true
}


func (i *MsgTransaction) DoSign(wallet wallets.IWallet) bool {
  if i.AddressFrom != wallet.GetAddress(i.CoinFrom) {
    return false
  }
  sign, ok := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), i.Hash())
  if !ok {
    return false
  }
  i.Sign = sign
  i.PubKey, ok = utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if !ok {
    return false
  }
  return true
}

func (i *MsgTransaction) DoVerify() bool {
  return utils.ECDSA256VerifySender(i.AddressFrom, i.PubKey, i.Hash(), i.Sign)
}
