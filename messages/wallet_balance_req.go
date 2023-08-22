package messages

import ( 
  "bytes"
  "encoding/gob"
  "crypto/sha512"
  
  "github.com/Lunkov/lib-wallets"
  
  "github.com/Lunkov/go-ecos-client/utils"
)

type GetBalanceReq struct {
  Address              string      `json:"address"      gorm:"column:address;primary_key"`
  Coin                 uint32      `json:"coin"`
  
  Sign                 []byte      `json:"sign"         gorm:"column:sign"`
  PublicKey            []byte      `json:"public_key"   gorm:"column:public_key"`
}

func NewGetBalanceReq() *GetBalanceReq {
  return &GetBalanceReq{}
}

func (m *GetBalanceReq) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(m)
  return buff.Bytes()
}

func (m *GetBalanceReq) Deserialize(msg []byte) error {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(m)
}

func (m *GetBalanceReq) Hash() []byte {
  sha := sha512.New()
  sha.Write([]byte(m.Address))
  sha.Write(utils.UInt32ToBytes(m.Coin))
  
  return sha.Sum(nil)
}

func (m *GetBalanceReq) Init(wallet wallets.IWallet, coin uint32) error {
  m.Address = wallet.GetAddress(coin)
  m.Coin = coin
  
  sign, err := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), m.Hash())
  if err != nil {
    return err
  }
  m.Sign = sign
  m.PublicKey, err = utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if err != nil {
    return err
  }
  return nil
}

func (m *GetBalanceReq) DoVerify() (bool, error) {
  return utils.ECDSA256VerifySender(m.Address, m.PublicKey, m.Hash(), m.Sign)
}
