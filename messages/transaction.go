package messages

import (
  "errors"
  "bytes"
  "time"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/gob"
  "encoding/json"
  "encoding/binary"
  
  "github.com/Lunkov/lib-wallets"
  
  "go-ecos-client/utils"
)

const (
  StatusTxNew       = 0
  StatusTxNewNFT    = 1
  StatusTxVerifying = 2
  StatusTxVerified  = 10
  StatusTxApproved  = 100
  StatusTxSaving    = 200
  StatusTxSaved     = 500
  StatusTxBadSign   = 10000
  StatusTxNoMoney   = 10001
)

type MsgTransaction struct {
  Version       uint32          `json:"version"`
  IdTx          uint32          `json:"id_tx"`

  AddressFrom   string          `json:"address_from"`
  AddressTo     string          `json:"address_to"`
  CoinFrom      uint32          `json:"coin_from"`
  CoinTo        uint32          `json:"coin_to"`
  ValueFrom     uint64          `json:"value"`
  ValueTo       uint64          `json:"value"`
  
  IdStatus      uint32          `json:"id_status"`
  
  UpdatedAt     time.Time       `json:"updated_at"`
  
  Sign          []byte          `json:"sign"         gorm:"column:sign"`
  PublicKey     []byte          `json:"public_key"   gorm:"column:public_key"`  
}

func NewMsgTransaction() *MsgTransaction {
  return &MsgTransaction{}
}


func (m *MsgTransaction) Init(idAction uint32,
                               wallet wallets.IWallet,
                               addressTo string, 
                               coinFrom uint32, 
                               coinTo uint32, 
                               valueFrom uint64,
                               valueTo uint64,
                               ) {
  m.AddressFrom = wallet.GetAddress(coinFrom)
  m.AddressTo = addressTo
  m.CoinFrom = coinFrom
  m.CoinTo = coinTo
  m.ValueFrom = valueFrom
  m.ValueTo = valueTo
  m.UpdatedAt = time.Now()
  m.IdTx = m.GenId()
}

func (m *MsgTransaction) InitNFT(idAction uint32,
                               wallet wallets.IWallet,
                               coin uint32,
                               ) {
  m.AddressFrom = wallet.GetAddress(coin)
  m.AddressTo = wallet.GetAddress(coin)
  m.CoinFrom = coin
  m.CoinTo = coin
  m.UpdatedAt = time.Now()
  m.IdTx = m.GenId()
}

func (m *MsgTransaction) GenId() uint32 {
  sha_256 := sha256.New()
  sha_256.Write(m.Hash())
  sha_256.Write([]byte(m.UpdatedAt.String()))
  return binary.LittleEndian.Uint32(sha_256.Sum(nil))
}

func (m *MsgTransaction) Hash() []byte {
  sha := sha512.New()
  sha.Write([]byte(m.AddressFrom))
  sha.Write([]byte(m.AddressTo))
  sha.Write(utils.UInt32ToBytes(m.CoinFrom))
  sha.Write(utils.UInt32ToBytes(m.CoinTo))
  sha.Write(utils.UInt64ToBytes(m.ValueFrom))
  sha.Write(utils.UInt64ToBytes(m.ValueTo))
  sha.Write(utils.UInt32ToBytes(m.IdStatus))
  
  return sha.Sum(nil)
}


func (m *MsgTransaction) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(m)
  return buff.Bytes()
}

func (m *MsgTransaction) Deserialize(msg []byte) error {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(m)
}

func (m *MsgTransaction) ToJSON() ([]byte, error) {
  jsonAnswer, err := json.Marshal(m)
  if err != nil {
    return nil, err
  }
  return jsonAnswer, nil
}

func (m *MsgTransaction) FromJSON(msg []byte) error {
  return json.Unmarshal(msg, m)
}

func (m *MsgTransaction) DoSign(wallet wallets.IWallet) error {
  if m.AddressFrom != wallet.GetAddress(m.CoinFrom) {
    return errors.New("DoSign: Wrong Address")
  }
  sign, err := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), m.Hash())
  if err != nil {
    return err
  }
  m.Sign = sign
  m.PublicKey, err = utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  return err
}

func (m *MsgTransaction) DoVerify() (bool, error) {
  return utils.ECDSA256VerifySender(m.AddressFrom, m.PublicKey, m.Hash(), m.Sign)
}
