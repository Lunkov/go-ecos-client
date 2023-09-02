package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestMsgOrganization(t *testing.T) {
  msg := NewOrganizationInfo()
  msg.CID = "12345"
  
  buf, errb := msg.Serialize()
  assert.Nil(t, errb) 
  
  msg2 := NewOrganizationInfo()
  
  assert.Nil(t, msg2.Deserialize(buf)) 
  
  assert.Equal(t, "12345", msg2.CID)
  assert.Equal(t, *msg, *msg2)
}
