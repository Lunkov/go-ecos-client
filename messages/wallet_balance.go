package messages

import ( 
  "sync"
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

func (m *GetBalanceReq) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(m)
  if err != nil {
    return false
  }
  return true
}

func (m *GetBalanceReq) Hash() []byte {
  sha := sha512.New()
  sha.Write([]byte(m.Address))
  sha.Write(utils.UInt32ToBytes(m.Coin))
  
  return sha.Sum(nil)
}

func (m *GetBalanceReq) Init(wallet wallets.IWallet, coin uint32) bool {
  m.Address = wallet.GetAddress(coin)
  m.Coin = coin
  
  sign, ok := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), m.Hash())
  if !ok {
    return false
  }
  m.Sign = sign
  m.PublicKey, ok = utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if !ok {
    return false
  }
  return true  
}

func (m *GetBalanceReq) DoVerify() bool {
  return utils.ECDSA256VerifySender(m.Address, m.PublicKey, m.Hash(), m.Sign)
}

type Balance struct {

  Address              string      `gorm:"column:address;primary_key"`
  Coin                 uint32      `gorm:"column:coin"`
  Balance              uint64      `gorm:"column:balance"`
  UnconfirmedBalance   uint64      `gorm:"column:unconfirmed_balance"`
  TotalReceived        uint64      `gorm:"column:total_received"`
  TotalSent            uint64      `gorm:"column:total_sent"`
  
  //LastTransaction      string      `gorm:"column:last_transaction"`
  //UpdatedAt            time.Time   `gorm:"column:updated_at;type:timestamp with time zone"`
  //Hash                 []byte      `gorm:"column:hash"`
}

func NewBalance() *Balance {
  return &Balance{}
}

func (b *Balance) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(b)
  return buff.Bytes()
}

func (b *Balance) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  err := decoder.Decode(b)
  if err != nil {
    return false
  }
  return true
}

type Balances struct {
  a     []*Balance
  mu    sync.RWMutex
}

func NewBalances() *Balances {
  return &Balances{a: make([]*Balance, 0)}
}

func (b *Balances) Add(i *Balance) {
  b.mu.Lock()
  b.a = append(b.a, i)
  b.mu.Unlock()
}

func (b *Balances) GetBalanses() ([]*Balance) {
  return b.a
}
