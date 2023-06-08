package objects

import (
  "bytes"
  "encoding/gob"
  "encoding/json"
)
/*
type WalletTxBalance struct {
  IdTx          []byte
  CIDBlock      string
  Timestamp     int64
  Debet         uint64
  Credit        uint64
}*/

type WalletTxDebet struct {
  IdTx          []byte
  CIDBlock      string
  Timestamp     int64
  Debet         uint64
  Credit        uint64
}

type WalletTxCredit struct {
  IdTxCredit    []byte
  IdTxDebet     []byte
  CIDBlock      string
  Timestamp     int64
  Credit        uint64
}

type WalletTxs struct {
  Address       string
  TxDebet     []WalletTxDebet
  TxCredit    []WalletTxCredit
}

func NewWalletTxEmpty() *WalletTxs {
  return &WalletTxs{TxDebet: make([]WalletTxDebet, 0), TxCredit: make([]WalletTxCredit, 0)}
}

func NewWalletTxBalance(address string) *WalletTxs {
  return &WalletTxs{Address: address, TxDebet: make([]WalletTxDebet, 0), TxCredit: make([]WalletTxCredit, 0)}
}

func (wtxb *WalletTxs) Serialize() ([]byte, bool) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(wtxb)
	if err != nil {
    return nil, false
	}

	return result.Bytes(), true
}

func (wtxb *WalletTxs) Deserialize(data []byte) bool {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(wtxb)
	if err != nil {
		return false
	}

	return true
}

func (wtxb *WalletTxs) ToJSON() ([]byte, bool) {
  jsonAnswer, err := json.Marshal(wtxb)
  if err != nil {
    return jsonAnswer, false
  }
	return jsonAnswer, true
}

func (wtxb *WalletTxs) FromJSON(data []byte) bool {
  if err := json.Unmarshal(data, wtxb); err != nil {
    return false
  }
  return true
}
