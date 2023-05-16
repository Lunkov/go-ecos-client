package engine

import (
  "time"
  "sync"
  "bytes"
  //"strconv"
  //"crypto/sha512"
  "crypto/sha256"
  //"github.com/jinzhu/gorm"
  _ "github.com/lib/pq"
  "encoding/gob"
  "github.com/golang/glog"
  //"github.com/google/uuid"
)

// https://github.com/hiromaily/go-crypto-wallet

const ChainID = uint64(202303)

// Block represents a block in the blockchain
type Block struct {
  ChainID          uint64
	Timestamp        int64
  
	Transactions     []*Transaction
	
  Hash             []byte
	Nonce            uint32
	Height           uint64
  
	PrevBlockHash    []byte
  PrevCID          string

  mu               sync.RWMutex
}

type BlockDB struct {
	Timestamp        int64
	CID              string
  PrevCID          string
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock() *Block {
	return &Block{
            ChainID: ChainID,
            Timestamp: time.Now().Unix(),
            PrevCID: "",
            PrevBlockHash: []byte{},
            Height: 0,
          }
}


// NewBlock creates and returns Block
func NewBlock(prevCID string, prevBlockHash []byte, nonce uint32, height uint64) *Block {
	return &Block{
            ChainID: ChainID,
            Timestamp: time.Now().Unix(),
            PrevCID: prevCID,
            PrevBlockHash: prevBlockHash,
            Nonce: nonce,
            Height: height,
          }
}


func (b *Block) AddTransaction(tx *Transaction) bool {
  b.mu.Lock()
  b.Transactions = append(b.Transactions, tx)
  b.mu.Unlock()
  return true
}

func (b *Block) HashBlock() []byte {
  h := sha256.New()
  h.Write(uint64ToBytes(b.ChainID))
  h.Write(int64ToBytes(b.Timestamp))
  h.Write(uint32ToBytes(b.Nonce))
  h.Write(uint64ToBytes(b.Height))
  h.Write([]byte(b.PrevCID))
  h.Write(b.PrevBlockHash)
  h.Write(b.HashTransactions())
  return h.Sum(nil)
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
  if len(b.Transactions) < 1 {
    return nil
  }
	var transactions [][]byte
	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)
	return mTree.RootNode.Data
}

// Serialize serializes the block
func (b *Block) Serialize() ([]byte, bool) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		glog.Errorf("ERR: Block.Serialize: %v", err)
    return nil, false
	}

	return result.Bytes(), true
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) (*Block, bool) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		glog.Errorf("ERR: Block.DeserializeBlock: %v", err)
    return nil, false
	}

	return &block, true
}

func (b *Block) SaveTransactionsToDB(db *DBO) (bool) {
  sql1 := db.handle.Table("transactions")
  tr := sql1.Begin()
  
  //if glog.V(7) {
    glog.Infof("LOG: Block.Height: %d transactions=%d", b.Height, len(b.Transactions))
  //}
  rdb := true
  for _, tx := range b.Transactions {
    rdb = tx.PutToDB(tr, b.Timestamp)
    if !rdb {
      break
    }
  }
  if rdb {
    tr.Commit()
  } else {
    tr.Rollback()
  }
  return rdb
}

func (b *Block) SaveBlockToDB(db *DBO, cid string) (bool) {
  sql1 := db.handle.Table("blocks")
  tr := sql1.Begin()
  bb := BlockDB{ Timestamp: b.Timestamp, CID: cid, PrevCID: b.PrevCID }
  err := tr.Create(&bb).Error
  if err == nil {
    tr.Commit()
  } else {
    glog.Errorf("ERR: Block.SaveBlockToDB: %v", err)
    tr.Rollback()
  }
  return err == nil
}

