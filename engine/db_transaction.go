package egine

import (
  //"fmt"
  "time"
  //"sync"
  //"strconv"
  "bytes"
  "encoding/gob"
  //"crypto/sha512"
  "crypto/sha256"
  //"crypto/ecdsa"
  "github.com/jinzhu/gorm"

  "github.com/golang/glog"
)

const version = byte(0x00)
// https://github.com/hiromaily/go-crypto-wallet
// https://en.bitcoin.it/wiki/Transaction

// TXInput represents a transaction input
type TXInput struct {
	Txid          []byte
	Vout          uint64
	Signature     []byte
	PubKey        []byte
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

type TransactionDB struct {
  Timestamp     int64       `gorm:"column:time_stamp"`
  Address       string      `gorm:"column:address;index:idx_address"`
  Vin           uint64      `gorm:"column:Vin"`
  Vout          uint64      `gorm:"column:Vout"`
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(address string, amount uint64) *Transaction {
	txin := TXInput{Vout: 0}
	txout := TXOutput{Value: amount, Address: address}
	tx := Transaction{ 
               Version: 1,
               Timestamp: time.Now().Unix(),
               Vin:  []TXInput{txin},
               Vout: []TXOutput{txout},
            }
	//tx.ID = tx.Hash()

	return &tx
}

// NewTX creates a new coinbase transaction
func NewTX(address string, amount uint64) *Transaction {
	txin := TXInput{Vout: 0}
	txout := TXOutput{Value: amount, Address: address}
	tx := Transaction{ 
               Version: 1,
               Timestamp: time.Now().Unix(),
               Vin:  []TXInput{txin},
               Vout: []TXOutput{txout},
            }
	//tx.ID = tx.Hash()

	return &tx
}

/*
// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock signs the output
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
*/
/*
// NewTXOutput create a new TXOutput
func NewTXOutput(value uint64, address string) *TXOutput {
	txo := &TXOutput{value, address}
	txo.Lock([]byte(address))

	return txo
}*/

// IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == 0
}

func (tx *Transaction) SetID() bool {
  var encoded bytes.Buffer
  var hash [32]byte

  encoder := gob.NewEncoder(&encoded)
  if err := encoder.Encode(tx); err != nil {
    glog.Errorf("ERR: Transaction.SetID: %v", err)
    return false
  }

  hash = sha256.Sum256(encoded.Bytes())
  tx.Id = hash[:]
  return true
}


// Serialize serializes the Transaction
func (i *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(i)
	if err != nil {
		glog.Errorf("ERR: Transaction.Serialize: %v", err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeTransaction(d []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&tx)
	if err != nil {
		glog.Errorf("ERR: Transaction.DeserializeTransaction: %v", err)
	}

	return &tx
}

func (tx *Transaction) PutToDB(tr *gorm.DB, timestamp int64) (bool) {
  for _, vin := range tx.Vin {
    if len(vin.PubKey) > 0 {
      //if glog.V(7) {
        glog.Infof("LOG: TX.Address.Minus: '%s' amount=%d", string(GetAddressPublicKey(vin.PubKey)), vin.Vout)
      //}
      tdb := TransactionDB{Timestamp: timestamp, Vout: vin.Vout, Address: string(GetAddressPublicKey(vin.PubKey))}
      
      err := tr.Create(&tdb).Error
      if err != nil {
        glog.Errorf("ERR: PutTxToDB() DB INSERT: %v", err)
        return false
      }
    }
  }
  for _, vout := range tx.Vout {
    //if glog.V(7) {
      glog.Infof("LOG: TX.Address.Plus: '%s' amount=%d", vout.Address, vout.Value)
    //}
    tdb := TransactionDB{Timestamp: timestamp, Vin: vout.Value, Address: vout.Address}
    err := tr.Create(&tdb).Error
    if err != nil {
      glog.Errorf("ERR: PutTxToDB() DB INSERT: %v", err)
      return false
    }
  }
  return true
}


/*

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			glog.Errorf("ERROR: Previous transaction is not correct")
      return false
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
		txCopy.Vin[inID].PubKey = nil
	}
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}


// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Verify verifies signatures of Transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.Txid)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Vin[inID].PubKey = nil
	}

	return true
}


// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(wallet *Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", wallet.GetAddress())
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}
*/
