package engine

import (
  "bytes"
  "time"
  "sync"
  "io/ioutil"
  //"crypto/sha512"
  ipfs "github.com/ipfs/go-ipfs-api"
  pubsub "github.com/libp2p/go-libp2p-pubsub"
  //"github.com/jinzhu/gorm"

  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-messages"
  "p2p-coins/api"
)

const genesisCID = "QmbMjXSexCPgykEGLQJzuDEA94cQiNTjSwLjezHzEAHTig"
const topicLastBlockCID = "last.block.cid"
const topicConsensus = "consensus"


type BlockPoolStat struct {
  NodeStatus            int32
  
  Height                uint64
  CountBlocks           uint64
  CountTransactions     uint64
  
  GenesisCID            string
  CurrentCID            string
  PrevCID               string
  
  TimeStart             time.Time
  TimeLastBlock         time.Time
  TimeStatusCons        time.Time
  
  MasterHost            string
  
  TimeSaveBlock         time.Duration
  TimeAddTransaction    time.Duration
  TimeSaveTransaction   time.Duration
  
  Nodes                 *NodesStatus
}

type BlockPool struct {
  CA                   *CAInfo

  DB                   *DBO
  IPFS                 *ipfs.Shell
  P2P                  *api.P2PConnect
  
  mu                    sync.RWMutex
  
  balances             *Balances
  
  oldBlock             *Block
  curBlock             *Block
  
  cfg                  *TokensCfg
  costTransactionSum    uint64
  
  cmdCh                 chan string
  
  stat                  BlockPoolStat
}

func NewBlockPool(ca *CAInfo, p2p *api.P2PConnect, db *DBO, ipfs *ipfs.Shell, cfg *TokensCfg) *BlockPool {
  bp := BlockPool{CA: ca, P2P: p2p, DB: db, IPFS: ipfs, cfg: cfg, balances: NewBalances(db)}
  bp.stat.NodeStatus = NodeSlave
  bp.stat.Nodes = NewNodesStatus()
  bp.costTransactionSum = 0
  for _, v := range cfg.Transactions.CostTransaction {
    bp.costTransactionSum += v.Value
  }
  if bp.P2P != nil {
    go bp.P2P.Subscribe(topicLastBlockCID, bp.SubRcvLastBlockCID)
    go bp.P2P.Subscribe(topicConsensus, bp.SubRcvConsensus)  
  }
  bp.cmdCh = make(chan string)
  bp.stat.TimeStart = time.Now()
  glog.Infof("LOG: Cost transaction: %d", bp.costTransactionSum)
  return &bp
}

func (b *BlockPool) GetLastCID() string {
  if b.curBlock == nil {
    return ""
  }
  return b.curBlock.PrevCID
}

func (b *BlockPool) GetStat() *BlockPoolStat {
  return &b.stat
}

func (b *BlockPool) GetLastBlock() Block {
  if b.oldBlock == nil {
    return Block{}
  }
  return (*b.oldBlock)
}

func (b *BlockPool) GetBalance(address string) (*messages.Balance, bool) {
  return b.balances.Get(address) 
}

func (b *BlockPool) PutToIPFS(data []byte) (string, bool) {
  // Add the data to IPFS
  cid, err := b.IPFS.Add(bytes.NewReader(data))
  if err != nil {
    glog.Errorf("ERR: putToIPFS: %v", err)
    return "", false
  }

  // Return the CID of the added data
  return cid, true
}

func (b *BlockPool) LoadLastBlock(lastCID string) bool {
  glog.Infof("LOG: Load Last Block: %s", lastCID)

  b.mu.Lock()
  defer b.mu.Unlock()

  var buf []byte
  var ok bool
  
  buf, ok = b.GetFromIPFS(lastCID)
  if !ok {
    return false
  }
  b.curBlock, ok = DeserializeBlock(buf)
  return ok
}

func (b *BlockPool) Reindex(initCID string, minCID string) bool {
  if initCID == "" {
    initCID = b.FindLastCID()
    glog.Infof("LOG: DB Find last cid=%s", initCID)
  }
  glog.Infof("LOG: Reindex from cid=%s", initCID)
  b.mu.Lock()
  defer b.mu.Unlock()

  var buf []byte
  var block *Block
  var ok bool
  cid := initCID
  
  for {
    glog.Infof("LOG: Reindex: Load cid=%s", cid)
    buf, ok = b.GetFromIPFS(cid)
    if !ok {
      return false
    }
    block, ok = DeserializeBlock(buf)
    if !ok {
      return false
    }
    if initCID == cid {
      b.stat.CurrentCID = cid
      b.stat.PrevCID = block.PrevCID
      b.stat.Height = block.Height
    }
    rdb := block.SaveTransactionsToDB(b.DB)   
    if !rdb {
      return false
    }
    rdb = block.SaveBlockToDB(b.DB, cid)
    if !rdb {
      return false
    }
    if block.PrevCID == "" || block.PrevCID == genesisCID {
      break
    }
    if block.PrevCID == minCID {
      break
    }
    cid = block.PrevCID
  }
  b.stat.CountTransactions = b.GetCountTransactions()
  b.stat.CountBlocks = b.GetCountBlocks()
  glog.Infof("LOG: STAT: Count Blocks=%d", b.stat.CountBlocks)
  glog.Infof("LOG: STAT: Count Transactions=%d", b.stat.CountTransactions)
  // check
  // 1. height --
  // 2. hash
  // 3. first is CID
  return true
}

/*
 * type TokenTransaction struct {
      AddressFrom        string  `json:"address_from"`
      AddressTo          string  `json:"address_to"`
      Coin               string  `json:"coin"`
      Value              uint64  `json:"value"`
      MaxCost            uint64  `json:"max_cost"`
      PublicKey          []byte  `json:"public_key"`
      Sign               []byte  `json:"sign"`
    }
 */
func (b *BlockPool) AddTransaction(tx *messages.TokenTransaction) bool {
  t1 := time.Now()
  //if glog.V(7) {
    glog.Infof("LOG: AddTransaction: '%s' -> '%s' = %d", tx.AddressFrom, tx.AddressTo, tx.Value)
  //}
  if tx.MaxCost < b.costTransactionSum {
    glog.Errorf("ERR: Cost transaction from '%s' is %d < %d cost transaction of node", tx.AddressFrom, tx.MaxCost, b.costTransactionSum)
    return false
  }
  if !PublicKeyECDSA256Verify(tx.PublicKey, tx.Hash(), tx.Sign) {
    glog.Errorf("ERR: PublicKeyECDSA256Verify: '%s' -> '%s' = %d", tx.AddressFrom, tx.AddressTo, tx.Value)
    return false
  }
  b.mu.Lock()
  balanceFrom, ok := b.balances.Get(tx.AddressFrom)
  if !ok {
    glog.Errorf("ERR: Balance '%s' not found", tx.AddressFrom)
    return false
  }
  if balanceFrom.Balance < (tx.Value + b.costTransactionSum) {
    glog.Errorf("ERR: Balance '%s' is smaller %d", tx.AddressFrom, (tx.Value + b.costTransactionSum))
    return false
  }
  ntx := NewTX(tx.AddressFrom, tx.Value + b.costTransactionSum)
  for _, v := range b.cfg.Transactions.CostTransaction {
    ntx.Vout = append(ntx.Vout, v)
  }
  ntx.SetID()
  b.curBlock.AddTransaction(ntx)
  b.stat.TimeAddTransaction = time.Since(t1)
  b.mu.Unlock()
  return true
}

func (b *BlockPool) SaveBlock() (bool) {
  b.mu.Lock()
  defer b.mu.Unlock()
  if b.curBlock == nil {
    return false
  }
  t1 := time.Now()
  buf, ok := b.curBlock.Serialize()
  if !ok {
    return false
  }
  newCID, okc := b.PutToIPFS(buf)
  if !okc {
    return false
  }
  glog.Infof("LOG: Last Block CID: %s", newCID)
  b.P2P.Publish(topicLastBlockCID, []byte(newCID))
  if !b.curBlock.SaveTransactionsToDB(b.DB) {
    return false
  }
  b.stat.CountTransactions += uint64(len(b.curBlock.Transactions))
  b.stat.CountBlocks += 1
  b.stat.Height = b.curBlock.Height + 1

  b.oldBlock = b.curBlock
  b.curBlock = NewBlock(newCID, b.curBlock.Hash, b.curBlock.Nonce, b.curBlock.Height + 1)
  b.stat.TimeSaveBlock = time.Since(t1)
  return true
}

func (b *BlockPool) GetFromIPFS(cid string) ([]byte, bool) {
  // Get the data associated with the CID
  if b.IPFS == nil {
    glog.Errorf("ERR: IPFS is not connected")
    return nil, false
  }
  reader, err := b.IPFS.Cat(cid)
  if err != nil {
    glog.Errorf("ERR: GetFromIPFS.Get: %v", err)
    return nil, false
  }

  // Read the data from the reader
  data, err := ioutil.ReadAll(reader)
  if err != nil {
    glog.Errorf("ERR: GetFromIPFS.Read: %v", err)
    return nil, false
  }

  // Return the data as a byte slice
  return data, true
}

func (b *BlockPool) SubRcvLastBlockCID(msg *pubsub.Message) {
  glog.Infof("LOG: SubRcvLastBlockCID: %v", msg)
  if string(msg.ReceivedFrom) == b.stat.MasterHost {
    // msg.Data
    // b.P2P.Publish(topicLastBlockCID, []byte(newCID))
    lastCID := string(msg.Data)
    b.Reindex(lastCID, b.stat.CurrentCID)
  }
}

func (b *BlockPool) SubRcvConsensus(msg *pubsub.Message) {
  glog.Infof("LOG: SubRcvConsensus: %v", msg)
  bn := NewNodesMsg()
  ok := bn.Deserialize(msg.Data)
  if !ok {
    glog.Infof("ERR: SubRcvConsensus.Deserialize: %v", msg)
    return
  }
  ok = b.CA.Verify(bn.Hash(), bn.Sign)
  if !ok {
    glog.Infof("ERR: SubRcvConsensus.Verify: %v", msg)
    return
  }
  // msg.ReceivedFrom == p2p.host.ID() {
}

func (b *BlockPool) Heartbit() {
  glog.Infof("LOG: Heartbit started")
  for {
 		select {
		case n := <-b.cmdCh:
			glog.Infof("LOG: Heartbit: CMD '%s'", n)
		case <-time.After(time.Millisecond * 5000):
      glog.Infof("LOG: Heartbit.After(time.Millisecond * 500)")
      //t1 := time.Now()
      if b.stat.NodeStatus == NodeMaster &&
        (b.curBlock == nil ||
        (b.curBlock != nil && (len(b.curBlock.Transactions) > 0 || time.Since(time.Unix(b.curBlock.Timestamp, 0)) > time.Second * 10))) {
			  glog.Infof("LOG: Heartbit: SaveBlock()")
        b.SaveBlock()
      }
      if b.stat.NodeStatus == NodeMaster {
        bn := NewNodesMsg()
        bn.CmdId = StatusCons
        bn.Data = b.stat.Nodes.Serialize()
        sign, ok := b.CA.Sign(bn.Hash())
        if !ok {
  			  glog.Errorf("ERR: Heartbit: topicConsensus().Sign")
        } else {
          bn.Sign = sign
	  		  glog.Infof("LOG: Heartbit: topicConsensus()")
          b.P2P.Publish(topicConsensus, bn.Serialize())
        }
      }
      glog.Infof("LOG: Heartbit.After(time.Millisecond * 1433)")
      if b.stat.NodeStatus == NodeSlave && b.stat.MasterHost == "" {
			  b.IamCandidat()
      }
      glog.Infof("LOG: Heartbit.After(time.Millisecond * 1733)")
      if b.stat.NodeStatus == NodeSlave && b.stat.MasterHost == "" {
			  b.IamMaster()
      }
    }
	}
}

func (b *BlockPool) IamCandidat() {
  glog.Infof("LOG: Heartbit: NO MASTER - StatusNodeCandidat")
  bn := NewNodesMsg()
  bn.CmdId = StatusNodeCandidat
  bnm := NodeStatus{NodeStatus: NodeMaster, HostAddress: b.P2P.GetHostID()}
  bn.Data = bnm.Serialize()
  sign, ok := b.CA.Sign(bn.Hash())
  if !ok {
    glog.Errorf("ERR: Heartbit: topicConsensus().Sign")
    return
  }
  bn.Sign = sign
  b.P2P.Publish(topicConsensus, bn.Serialize())
}

func (b *BlockPool) IamMaster() {
  glog.Infof("LOG: Heartbit: NO MASTER")
  bn := NewNodesMsg()
  bn.CmdId = StatusNode
  bnm := NodeStatus{NodeStatus: NodeMaster, HostAddress: b.P2P.GetHostID()}
  bn.Data = bnm.Serialize()
  sign, ok := b.CA.Sign(bn.Hash())
  if !ok {
    glog.Errorf("ERR: Heartbit: topicConsensus().Sign")
    return
  }
  bn.Sign = sign
  b.P2P.Publish(topicConsensus, bn.Serialize())
  b.stat.NodeStatus = NodeMaster
  b.stat.MasterHost = b.P2P.GetHostID()
}

// CountTransactions returns the number of transactions
func (b *BlockPool) GetCountTransactions() uint64 {
	var count uint64
  b.DB.handle.Table("transactions").Count(&count)
	return count
}

// GetCountBlocks returns the number of blocks
func (b *BlockPool) GetCountBlocks() uint64 {
	var count uint64
  b.DB.handle.Table("blocks").Count(&count)
	return count
}

func (b *BlockPool) FindLastCID() string {
  var block BlockDB
	b.DB.handle.Table("blocks").Order("timestamp desc").FirstOrInit(&block)
	return block.CID
}
