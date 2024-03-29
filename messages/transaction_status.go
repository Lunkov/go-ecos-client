package messages

import (
  "bytes"
  "time"
  "crypto/sha512"
  "encoding/gob"
  
  "github.com/Lunkov/lib-wallets"

  "github.com/Lunkov/go-ecos-client/utils"
  "github.com/Lunkov/go-ecos-client/objects"
)

type MsgTransactionStatus struct {
  Version       uint32          `json:"version"`
  IdTx          []byte          `json:"id_tx"`
  IdStatus      uint32          `json:"id_status"`
  Timestamp     int64
  
  Sign          []byte          `json:"sign"         gorm:"column:sign"`
  PublicKey     []byte          `json:"public_key"   gorm:"column:public_key"`  
}

func NewMsgTransactionStatus(IdTx []byte) *MsgTransactionStatus {
  return &MsgTransactionStatus{Version: objects.TXVersion, IdTx: IdTx, Timestamp: time.Now().Unix()}
}

func (m *MsgTransactionStatus) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(m)
  return buff.Bytes()
}

func (m *MsgTransactionStatus) Deserialize(msg []byte) error {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(m)
}


func (m *MsgTransactionStatus) Hash512() []byte {
  sha := sha512.New()
  sha.Write(utils.UInt32ToBytes(m.Version))
  sha.Write(m.IdTx)
  sha.Write(m.PublicKey)
  
  return sha.Sum(nil)
}

func (m *MsgTransactionStatus) DoSign(wallet wallets.IWallet) error {
  var err error
  m.PublicKey, err = utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if err != nil {
    return err
  }
  sign, errs := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), m.Hash512())
  if errs != nil {
    return errs
  }
  m.Sign = sign
  return nil
}

func (m *MsgTransactionStatus) DoVerify() (bool, error) {
  return utils.ECDSA256VerifyHash512(m.PublicKey, m.Hash512(), m.Sign)
}
