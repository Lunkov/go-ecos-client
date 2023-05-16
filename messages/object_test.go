package messages

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestMessageObject(t *testing.T) {
	
  obj := ReqActionObject()
  obj.Init("token", "user", "login@domain.org", "admin")
  
  dbui.Hash()
     
  uio := NewReqUserInviteFromOrg()
  uio.Init("ANO", dbui)
  assert.Equal(t, 64, len(uio.Hash()))
  
  msg := uio.Pack()
  assert.Equal(t, 284, len(msg))
  
  uio2 := NewReqUserInviteFromOrg()
  uio2.Unpack(msg)
  
  assert.Equal(t, uio.NodeUrl, uio2.NodeUrl)
  assert.Equal(t, uio.Login, uio2.Login)
  assert.Equal(t, uio.EMail, uio2.EMail)
}

func BenchmarkMutexMap(b *testing.B) {

  b.ResetTimer()

  dbui := ReqActionObject()
  dbui.Init("token", "user", "login@domain.org", "admin")
  
  dbui.Hash()
     
  uio := NewReqUserInviteFromOrg()
  uio.Init("ANO", dbui)
  assert.Equal(b, 64, len(uio.Hash()))
  
  msg := uio.Pack()
  assert.Equal(b, 244, len(msg))
  
  uio2 := NewReqUserInviteFromOrg()
  uio2.Unpack(msg)
  
  assert.Equal(b, uio.NodeUrl, uio2.NodeUrl)
  assert.Equal(b, uio.Login, uio2.Login)
  assert.Equal(b, uio.EMail, uio2.EMail)
}
