package objects

import (
  "errors"
  "bytes"
  "time"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/gob"
  
  "github.com/Lunkov/go-hdwallet"
  "github.com/Lunkov/lib-wallets"
  
  "go-ecos-client/utils"
)

const TXVersion = uint32(0x00)
// https://github.com/hiromaily/go-crypto-wallet
// https://en.bitcoin.it/wiki/Transaction

// TXInput represents a transaction input
type TXInput struct {
  Txid          []byte
  Address       string
  CIDNFT        string
  Vout          uint64
  Signature     []byte
  PublicKey     []byte
}


// TXOutput represents a transaction output
type TXOutput struct {
  Address        string                 `yaml:"address"`
  CIDNFT         string                 `yaml:"cid"`
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
    return nil, false
  }

  hash = sha256.Sum256(encoded.Bytes())
  return hash[:], true
}

// Serialize serializes the Transaction
func (tx *Transaction) Serialize() ([]byte, error) {
  var result bytes.Buffer
  encoder := gob.NewEncoder(&result)

  err := encoder.Encode(tx)
  if err != nil {
    return nil, err
  }

  return result.Bytes(), nil
}

// DeserializeBlock deserializes a block
func (tx *Transaction) Deserialize(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(tx)
	if err != nil {
    return err
	}

	return nil
}

func (m TXInput) Hash512(timestamp int64) []byte {
  sha := sha512.New()
  sha.Write(m.Txid)
  sha.Write([]byte(m.Address))
  sha.Write(utils.Int64ToBytes(timestamp))
  sha.Write(utils.UInt64ToBytes(m.Vout))
  sha.Write(m.PublicKey)
  
  return sha.Sum(nil)
}

func (tx *Transaction) GetValueIn() uint64 {
  vin := uint64(0)
  for _, v := range tx.Vin {
    vin += v.Vout
  }
  return vin
}

func (tx *Transaction) GetValueOut() uint64 {
  vout := uint64(0)
  for _, v := range tx.Vout {
    vout += v.Value
  }
  return vout
}


func (tx *Transaction) DoSign(wallet wallets.IWallet) error {
  address := wallet.GetAddress(hdwallet.ECOS)
  PublicKey, err := utils.ECDSAPublicKeySerialize(wallet.GetECDSAPublicKey())
  if err != nil {
    return err
  }
  for inID, vin := range tx.Vin {
    if vin.Txid == nil {
      return errors.New("DoSign: vin.Txid is empty")
    }
    tx.Vin[inID].Address = address
    tx.Vin[inID].PublicKey = PublicKey
    sign, errs := utils.ECDSA256SignHash512(wallet.GetECDSAPrivateKey(), tx.Vin[inID].Hash512(tx.Timestamp))
    if errs != nil {
      return errs
    }
    tx.Vin[inID].Signature = sign
  }
  return nil
}

func (tx *Transaction) DoVerify() (bool, error) {
  for inID, vin := range tx.Vin {
    if vin.Txid == nil {
      return false, errors.New("DoVerify: vin.Txid is empty")
    }
    vok, errv := utils.ECDSA256VerifyHash512(tx.Vin[inID].PublicKey, tx.Vin[inID].Hash512(tx.Timestamp), tx.Vin[inID].Signature)
    if !vok || errv != nil {
      return false, errv
    }
  }
  return true, nil
}
