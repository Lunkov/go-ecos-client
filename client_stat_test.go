package client

import (
  "flag"
  //"time"
  "testing"
  "github.com/stretchr/testify/assert"
  //"github.com/Lunkov/go-ecos-client/messages"

)

func TestStat(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  //flag.Set("v", "9")
  //flag.Parse()

  client := NewClientECOS([]string{"http://127.0.0.1:8085/"}, 3)
  
  stat1, err := client.GetNetworkStat()
  assert.Nil(t, err)

  if stat1 != nil {
    //assert.Equal(t, messages.NetworkStat{NodeStatus:0, Height:0x0, CountBlocks:0x0, CountTransactions:0x0, GenesisCID:"", CurrentCID:"", PrevCID:"", TimeStart:time.Date(2023, time.June, 11, 9, 47, 16, 764230294, time.Local), TimeLastBlock:time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), TimeStatusCons:time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), MasterHost:"", TimeSaveBlock:0, TimeAddTransaction:0, TimeSaveTransaction:0}, *stat1)
  }

  stat2, err2 := client.GetNetworkStat()
  assert.Nil(t, err2)
  
  if stat2 != nil {
    //assert.Equal(t, messages.NetworkStat{NodeStatus:0, Height:0x0, CountBlocks:0x0, CountTransactions:0x0, GenesisCID:"", CurrentCID:"", PrevCID:"", TimeStart:time.Date(2023, time.June, 11, 9, 47, 16, 764230294, time.Local), TimeLastBlock:time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), TimeStatusCons:time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), MasterHost:"", TimeSaveBlock:0, TimeAddTransaction:0, TimeSaveTransaction:0}, *stat2)
  }
}
