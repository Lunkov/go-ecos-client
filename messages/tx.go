package messages

import (
  "bytes"
  "time"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/gob"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-wallets"
  "github.com/Lunkov/go-ecos-client/utils"
)

const TXVersion = uint32(0x00)
// https://github.com/hiromaily/go-crypto-wallet
// https://en.bitcoin.it/wiki/Transaction

// TXInput represents a transaction input
type TXInput struct {
	Txid          []byte
	Vout          uint64
	Signature     []byte
	PublicKey     []byte
}


// TXOutput represents a transaction output
type TXOutput struct {
	Address        string                 `yaml:"address"`
	Value          uint64                 `yaml:"value"`
}


type Transaction struct {
  Id            []byte
  Version       uint32
	Timestamp     int64
  
	Vin      []TXInput
	Vout     []TXOutput
}

func NewTX() *Transaction {
  return &Transaction{ 
               Version: TXVersion,
               Timestamp: time.Now().Unix(),
               Vin:  make([]TXInput, 0),
               Vout: make([]TXOutput, 0),
            }
}

// IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == 0
}

func (tx *Transaction) CalcID() ([]byte, bool) {
  var encoded bytes.Buffer
  var hash [32]byte

  encoder := gob.NewEncoder(&encoded)
  if err := encoder.Encode(tx); err != nil {
    glog.Errorf("ERR: Transaction.SetID: %v", err)
    return nil, false
  }

  hash = sha256.Sum256(encoded.Bytes())
  return hash[:], true
}

// Serialize serializes the Transaction
func (tx *Transaction) Serialize() ([]byte, bool) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		glog.Errorf("ERR: Transaction.Serialize: %v", err)
    return nil, false
	}

	return result.Bytes(), true
}

// DeserializeBlock deserializes a block
func (tx *Transaction) Deserialize(data []byte) bool {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(tx)
	if err != nil {
		glog.Errorf("ERR: Transaction.DeserializeTransaction: %v", err)
    return false
	}

	return true
}
  
func (m TXInput) Hash512(timestamp int64) []byte {
  sha := sha512.New()
  sha.Write(m.Txid)
  sha.Write(utils.Int64ToBytes(timestamp))
  sha.Write(utils.UInt64ToBytes(m.Vout))
  sha.Write(m.PublicKey)
  
  return sha.Sum(nil)
}

func (tx *Transaction) DoSign(wallet wallets.IWallet) bool {
  PublicKey, ok := utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if !ok {
    return false
  }
  for inID, vin := range tx.Vin {
    if vin.Txid == nil {
      glog.Errorf("ERR: Previous transaction is not correct")
      return false
    }
    tx.Vin[inID].PublicKey = PublicKey
    sign, ok := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), tx.Vin[inID].Hash512(tx.Timestamp))
    if !ok {
      return false
    }
    tx.Vin[inID].Signature = sign
  }
  return true
}

func (tx *Transaction) DoVerify() bool {
  for inID, vin := range tx.Vin {
    if vin.Txid == nil {
      glog.Errorf("ERR: Previous transaction is not correct")
      return false
    }
    if !utils.ECDSA256VerifyHash512(tx.Vin[inID].PublicKey, tx.Vin[inID].Hash512(tx.Timestamp), tx.Vin[inID].Signature) {
      return false
    }
  }
  return true
}
