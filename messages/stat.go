package messages

import (
  "time"
  "bytes"
  "encoding/gob"
)

type NetworkStat struct {
  NodeStatus            int32
  
  Height                uint64
  CountBlocks           uint64
  CountTransactions     uint64
  SumTransactions       uint64
  CountWallets          uint64
  
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
  
}

func NewNetworkStat() *NetworkStat {
  return &NetworkStat{}
}

func (i *NetworkStat) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *NetworkStat) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(i) == nil
}


type NodeStat struct {
  NodeStatus            int32
  
  Height                uint64
  CountBlocks           uint64
  CountTransactions     uint64
  SumTransactions       uint64
  CountWallets          uint64
  
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
  
}

func NewNodeStat() *NodeStat {
  return &NodeStat{}
}

func (i *NodeStat) Serialize() []byte {
  var buff bytes.Buffer
  encoder := gob.NewEncoder(&buff)
  encoder.Encode(i)
  return buff.Bytes()
}

func (i *NodeStat) Deserialize(msg []byte) bool {
  buf := bytes.NewBuffer(msg)
  decoder := gob.NewDecoder(buf)
  return decoder.Decode(i) == nil
}
