package messages

import (
  "time"
  "bytes"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/gob"
  "encoding/binary"
  
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
  "github.com/Lunkov/go-ecos-client/utils"
)

type ReqActionObject struct {
  Version       string          `json:"version"`
  
  MessageId     uint32          `json:"message_id"`
  
  ObjectId      string          `json:"object_id"`

  ActionId      uint32          `json:"action_id"`

  CID           string          `json:"cid"`

  DataObject    []byte          `json:"data_object"`

  UpdatedAt     time.Time       `json:"updated_at"`

  Address       string          `json:"address"`

  PublicKey     []byte          `json:"pubkey"`
  Sign          []byte          `json:"sign"`
}

func NewReqActionObject() *ReqActionObject {
  return &ReqActionObject{Version: "1"}
}

func (i *ReqActionObject) Init(idAction uint32, idObject string, dataObject []byte) {
  i.ActionId = idAction
  i.ObjectId = idObject
  i.DataObject = dataObject
  i.UpdatedAt = time.Now()
  i.MessageId = i.msgId()
}

func (i *ReqActionObject) msgId() uint32 {
  sha_256 := sha256.New()
  sha_256.Write([]byte(i.ObjectId + i.CID))
  sha_256.Write(utils.UInt32ToBytes(i.ActionId))
  sha_256.Write(i.DataObject)
  sha_256.Write([]byte(i.UpdatedAt.String()))
  return binary.LittleEndian.Uint32(sha_256.Sum(nil))
}

func (i *ReqActionObject) Hash() []byte {
  sha_512 := sha512.New()
  sha_512.Write([]byte(i.Version + i.ObjectId + i.CID))
  sha_512.Write(utils.UInt32ToBytes(i.ActionId))
  sha_512.Write(utils.UInt32ToBytes(i.MessageId))
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

func (i *ReqActionObject) DoSign(wallet wallets.IWallet) bool {
  sign, ok := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), i.Hash())
  if !ok {
    return false
  }
  i.Address = wallet.GetAddress(hdwallet.ECOS)
  i.Sign = sign
  i.PublicKey, ok = utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if !ok {
    return false
  }
  return true
}

func (i *ReqActionObject) DoVerify() bool {
  return utils.ECDSA256VerifySender(i.Address, i.PublicKey, i.Hash(), i.Sign)
}
