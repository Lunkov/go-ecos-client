package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestMsgPassport(t *testing.T) {
  msg := NewPassportInfo()
  msg.CID = "12345"
  
  buf, okb := msg.Serialize()
  assert.True(t, okb) 
  
  msg2 := NewPassportInfo()
  
  assert.True(t, msg2.Deserialize(buf)) 
  
  assert.Equal(t, "12345", msg2.CID)
  assert.Equal(t, *msg, *msg2)
}
