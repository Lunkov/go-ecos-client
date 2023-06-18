package objects

import (
  "bytes"
  "sort"
  "encoding/gob"
  "encoding/json"
)

const (
  DirectionUndef  = 0
  DirectionGen    = 1
  DirectionInput  = 2
  DirectionOutput = 3
)

type WalletTransaction struct {
  IdTx          []byte
  CIDBlock        string
  Timestamp       int64

  DirectionCoins  uint32  
  Amount          uint64

  DirectionNFT    uint32  
  CIDNFT          string

  Balance         uint64
}

type WalletTransactions struct {
  Address       string
  Transactions  []WalletTransaction
}

func NewWalletTransactionsEmpty() *WalletTransactions {
  return &WalletTransactions{Transactions: make([]WalletTransaction, 0)}
}

func NewWalletTransactionsBalance(address string) *WalletTransactions {
  return &WalletTransactions{Address: address, Transactions: make([]WalletTransaction, 0)}
}

func (wtxb *WalletTransactions) SortByDate() {
  sort.Slice(wtxb.Transactions, func(i, j int) bool {
    return wtxb.Transactions[i].Timestamp < wtxb.Transactions[j].Timestamp
  })
}

func (wtxb *WalletTransactions) SortByDateDesc() {
  sort.Slice(wtxb.Transactions, func(i, j int) bool {
    return wtxb.Transactions[i].Timestamp > wtxb.Transactions[j].Timestamp
  })
}

func (wtxb *WalletTransactions) FilterNFT(nft string) *WalletTransactions {
  nwtx := NewWalletTransactionsBalance(wtxb.Address)
  for _, v := range wtxb.Transactions {
    if v.CIDNFT == nft {
      nwtx.Transactions = append(nwtx.Transactions, v)
    }
  } 
  sort.Slice(nwtx.Transactions, func(i, j int) bool {
    return nwtx.Transactions[i].Timestamp < nwtx.Transactions[j].Timestamp
  })
  return nwtx
}

func (wtxb *WalletTransactions) Have(nft string) bool {
  have := false
  for _, v := range wtxb.Transactions {
    if v.CIDNFT == nft {
      have = (v.DirectionNFT == DirectionInput || v.DirectionNFT == DirectionGen) 
    }
  } 
  return have
}

func (wtxb *WalletTransactions) Append(wt *WalletTransactions) {
  if wtxb.Address == wt.Address {
    for _, v := range wt.Transactions {
      wtxb.Transactions = append(wtxb.Transactions, v)
    }
  }
}

func (wtxb *WalletTransactions) GetStat() (uint64, uint64) {
  totalReceived := uint64(0)
  totalSent := uint64(0)

  for _, v := range wtxb.Transactions {
    if v.DirectionCoins == DirectionInput || v.DirectionCoins == DirectionGen {
      totalReceived += v.Amount
      continue
    }
    if v.DirectionCoins == DirectionOutput {
      totalSent += v.Amount
    }
  } 
  return totalReceived, totalSent
}

func (wtxb *WalletTransactions) RecalcBalance(startBalance uint64) {
  wtxb.SortByDate()
  
  balance := startBalance
  for k, v := range wtxb.Transactions {
    if v.DirectionCoins == DirectionInput || v.DirectionCoins == DirectionGen {
      balance += v.Amount
    }
    if v.DirectionCoins == DirectionOutput {
      balance -= v.Amount
    }
    wtxb.Transactions[k].Balance = balance
  } 
}

func (wtxb *WalletTransactions) Serialize() ([]byte, bool) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(wtxb)
	if err != nil {
    return nil, false
	}

	return result.Bytes(), true
}

func (wtxb *WalletTransactions) Deserialize(data []byte) bool {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(wtxb)
	if err != nil {
		return false
	}

	return true
}

func (wtxb *WalletTransactions) ToJSON() ([]byte, bool) {
  jsonAnswer, err := json.Marshal(wtxb)
  if err != nil {
    return jsonAnswer, false
  }
	return jsonAnswer, true
}

func (wtxb *WalletTransactions) FromJSON(data []byte) bool {
  if err := json.Unmarshal(data, wtxb); err != nil {
    return false
  }
  return true
}
