package objects

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/Lunkov/go-ecos-client/utils"
)

func TestMsgPassport(t *testing.T) {
  key, ok := utils.RSAGenerate(2048)
  assert.True(t, ok)
  
  msg := NewPassport()
  msg.DisplayName = "Test User"
  
  msg.IssuerPubKey, ok = utils.RSASerializePublicKey(&key.PublicKey)
  assert.True(t, ok)
  
  msg.IssuerSign, ok = utils.RSASign(key, msg.Hash())
  assert.True(t, ok)
  
  buf, okb := msg.Serialize()
  assert.True(t, okb) 
  
  msg2 := NewPassport()
  
  assert.True(t, msg2.Deserialize(buf)) 
  
  assert.True(t, utils.RSAVerify(msg.IssuerPubKey, msg.Hash(), msg.IssuerSign))
  
  assert.Equal(t, "Test User", msg2.DisplayName)
  assert.Equal(t, *msg, *msg2)
}
